package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"github.com/pierrec/lz4"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

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

func (h *uploadHandlersV3) StartWorker(ctx context.Context, index int) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.socketServerAddr)
	if err != nil {
		return err
	}
	h.conns[index] = conn
	go func() {
		// TODO: may block here
		<-ctx.Done()
		if h.files[index] != nil {
			h.files[index].Close()
		}
	}()
	return nil
}

func (h *uploadHandlersV3) reconnect(ctx context.Context, index int) error {
	h.conns[index].Close()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.socketServerAddr)
	if err != nil {
		return err
	}
	h.conns[index] = conn
	return nil
}

func (h *uploadHandlersV3) OnFileError(filePath string, err error) {
	h.uploadProcess.OnFileError(filePath, err)
}

func (h *uploadHandlersV3) formatPath(filePath string) ([]string, error) {
	rel, err := filepath.Rel(h.srcPath, filePath)
	if err != nil {
		return nil, err
	}
	pathList := strings.Split(filepath.Join(h.dstPath, rel), string(os.PathSeparator))
	return pathList, nil
}

func (h *uploadHandlersV3) DirHandler(ctx context.Context, filePath string, infos []os.FileInfo, continues []bool) error {
	dirPath, err := h.formatPath(filePath)
	if err != nil {
		return err
	}
	uploadReqDirItemV3 := make([]*pb.UploadReqDirItemV3, 0, cap(infos))
	for _, info := range infos {
		modifyTime := uint64(info.ModTime().UnixNano())
		uploadReqDirItemV3 = append(uploadReqDirItemV3, &pb.UploadReqDirItemV3{
			Name:       info.Name(),
			Mode:       uint64(info.Mode()),
			Size:       uint64(info.Size()),
			CreateTime: modifyTime,
			ModifyTime: modifyTime,
			ChangeTime: modifyTime,
			AccessTime: modifyTime,
		})
	}
	_, err = ReqRespWithConn(h.conn, rpcutil.CommandUploadV2Dir, &pb.UploadReqV3{
		DriverId:           h.driverId,
		DirPath:            dirPath,
		UploadReqDirItemV3: uploadReqDirItemV3,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (h *uploadHandlersV3) copyFile(conn net.Conn, f *os.File, size int64) error {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "CopyStart", size)
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
	println(conn.RemoteAddr().String(), "CopyEnd", n)
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

func (h *uploadHandlersV3) uploadFile(ctx context.Context, conn net.Conn, index int, filePath string, info os.FileInfo, dirPath []string, name string) (exist bool, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	h.files[index] = f
	defer func() {
		h.files[index] = nil
		f.Close()
	}()
	size := uint64(info.Size())
	hash, err := h.getHash(f)
	if err != nil {
		return
	}

	println(conn.RemoteAddr().String(), "uploadFile", filePath, hash)

	modifyTime := uint64(info.ModTime().UnixNano())
	status, err := ReqRespWithConn(conn, rpcutil.CommandUploadV2File, &pb.UploadReqV2{
		DriverId:   h.driverId,
		DirPath:    dirPath,
		Name:       name,
		Hash:       hash,
		Mode:       uint64(info.Mode()),
		Size:       size,
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	}, nil)
	if err != nil {
		return
	}
	if status == rpcutil.ENotExist {
		println(conn.RemoteAddr().String(), "encoder", len(h.encoder), h.encoder)
		err = rpcutil.WriteString(conn, h.encoder)
		if err != nil {
			return
		}

		err = h.copyFile(conn, f, info.Size())
		if err != nil {
			return
		}

		_, _, err = rpcutil.ReadStatus(conn)
		if err != nil {
			return
		}
		return
	}
	exist = true

	return
}
