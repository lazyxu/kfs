package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"github.com/pierrec/lz4"

	"github.com/lazyxu/kfs/pb"
)

type uploadHandlersV3 struct {
	core.DefaultWalkDirHandlers
	uploadProcess    core.UploadDirProcess
	concurrent       int
	encoder          string
	verbose          bool
	socketServerAddr string
	conns            []net.Conn
	files            []*os.File
	driverId         uint64
	srcPath          string
	dstPath          string
	conn             net.Conn
}

func (h *uploadHandlersV3) OnFileError(filePath string, err error) {
	h.uploadProcess.OnFileError(filePath, err)
}

func (h *uploadHandlersV3) formatPath(filePath string) ([]string, error) {
	rel, err := filepath.Rel(h.srcPath, filePath)
	if err != nil {
		return nil, err
	}
	actualPath := filepath.Join(h.dstPath, rel)
	pathList := strings.Split(actualPath, string(os.PathSeparator))
	newPathList := []string{}
	for _, path := range pathList {
		if path != "" {
			newPathList = append(newPathList, path)
		}
	}
	return newPathList, nil
}

func (h *uploadHandlersV3) DirHandler(ctx context.Context, filePath string, dirInfo os.FileInfo, infos []os.FileInfo, continues []bool) error {
	dirPath, err := h.formatPath(filePath)
	if err != nil {
		return err
	}
	uploadReqDirItemCheckV3 := make([]*pb.UploadReqDirItemCheckV3, cap(infos))
	for i, info := range infos {
		h.uploadProcess.PushFile(info)
		modifyTime := uint64(info.ModTime().UnixNano())
		uploadReqDirItemCheckV3[i] = &pb.UploadReqDirItemCheckV3{
			Name:       info.Name(),
			Size:       uint64(info.Size()),
			ModifyTime: modifyTime,
		}
	}

	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	var respCheck pb.UploadRespV3
	_, err = ReqRespWithConn(h.conn, rpcutil.CommandUploadV3DirCheck, &pb.UploadReqCheckV3{
		DriverId:                h.driverId,
		DirPath:                 dirPath,
		UploadReqDirItemCheckV3: uploadReqDirItemCheckV3,
	}, &respCheck)
	if err != nil {
		return err
	}

	uploadReqDirItemV3 := make([]*pb.UploadReqDirItemV3, cap(infos))
	for i, hash := range respCheck.Hash {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
		info := infos[i]
		p := filepath.Join(filePath, info.Name())
		h.uploadProcess.StartFile(p, info)
		if !info.IsDir() && hash == "" {
			hash, err = h.uploadFile(ctx, p, info)
			if err != nil {
				return err
			}
		}
		modifyTime := uint64(info.ModTime().UnixNano())
		uploadReqDirItemV3[i] = &pb.UploadReqDirItemV3{
			Name:       info.Name(),
			Hash:       hash,
			Mode:       uint64(info.Mode()),
			Size:       uint64(info.Size()),
			CreateTime: modifyTime,
			ModifyTime: modifyTime,
			ChangeTime: modifyTime,
			AccessTime: modifyTime,
		}
		h.uploadProcess.EndFile(p, info)
	}

	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	h.uploadProcess.StartDir(filePath, uint64(len(uploadReqDirItemV3)))
	_, err = ReqRespWithConn(h.conn, rpcutil.CommandUploadV3Dir, &pb.UploadReqV3{
		DriverId:           h.driverId,
		DirPath:            dirPath,
		UploadReqDirItemV3: uploadReqDirItemV3,
	}, nil)
	if err != nil {
		return err
	}
	h.uploadProcess.EndDir(filePath, dirInfo)

	return nil
}

func (h *uploadHandlersV3) copyFile(conn net.Conn, f *os.File, size int64) error {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	println("CopyStart", size)
	var n int64
	if h.encoder == "lz4" {
		w := lz4.NewWriter(conn)
		n, err = io.CopyN(w, f, size)
		if err != nil {
			return err
		}
		defer w.Flush()
	} else {
		n, err = io.CopyN(conn, f, size)
		if err != nil {
			return err
		}
	}
	println("CopyEnd", n)
	return nil
}

func (h *uploadHandlersV3) getHash(f *os.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (h *uploadHandlersV3) uploadFile(ctx context.Context, filePath string, info os.FileInfo) (hash string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()
	size := uint64(info.Size())
	hash, err = h.getHash(f)
	if err != nil {
		return
	}

	select {
	case <-ctx.Done():
		err = context.Canceled
		return
	default:
	}
	println("uploadFile", filePath, hash)

	status, err := ReqRespWithConn(h.conn, rpcutil.CommandUploadV3File, &pb.UploadFileV3{
		Hash: hash,
		Size: size,
	}, nil)
	if err != nil {
		return
	}
	if status == rpcutil.ENotExist {
		select {
		case <-ctx.Done():
			err = context.Canceled
			return
		default:
		}
		println("encoder", len(h.encoder), h.encoder)
		err = rpcutil.WriteString(h.conn, h.encoder)
		if err != nil {
			return
		}

		err = h.copyFile(h.conn, f, info.Size())
		if err != nil {
			return
		}

		_, _, err = rpcutil.ReadStatus(h.conn)
		if err != nil {
			return
		}
		return
	}
	return
}
