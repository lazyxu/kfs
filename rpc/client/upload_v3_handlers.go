package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dustin/go-humanize"
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
	uploadDeviceId   string
	uploadTime       uint64
}

func (h *uploadHandlersV3) FilePathFilter(filePath string) bool {
	return h.uploadProcess.FilePathFilter(filePath)
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
	dirPath, err := h.formatPath(filepath.Dir(filePath))
	if err != nil {
		h.uploadProcess.OnFileError(filePath, err)
		return nil
	}
	uploadReqDirItemCheckV3 := make([]*pb.UploadReqDirItemCheckV3, len(infos))
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

	var startResp pb.UploadStartDirResp
	dirModifyTime := uint64(dirInfo.ModTime().UnixNano())
	_, err = ReqRespWithConn(h.conn, rpcutil.CommandUploadStartDir, &pb.UploadStartDirReq{
		DriverId:                h.driverId,
		DirPath:                 dirPath,
		Name:                    dirInfo.Name(),
		Hash:                    "",
		Mode:                    uint64(dirInfo.Mode()),
		Size:                    uint64(dirInfo.Size()),
		CreateTime:              dirModifyTime,
		ModifyTime:              dirModifyTime,
		ChangeTime:              dirModifyTime,
		AccessTime:              dirModifyTime,
		UploadDeviceId:          h.uploadDeviceId,
		UploadTime:              h.uploadTime,
		UploadReqDirItemCheckV3: uploadReqDirItemCheckV3,
	}, &startResp)
	if err != nil {
		return err
	}

	for i, hash := range startResp.Hash {
		info := infos[i]
		p := filepath.Join(filePath, info.Name())
		if !info.IsDir() {
			h.uploadProcess.StartFile(p, info)
		}
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
		if !info.IsDir() {
			if hash == "" {
				var fileErr error
				hash, fileErr, err = h.uploadFile(ctx, p, info)
				if fileErr != nil {
					h.uploadProcess.OnFileError(p, fileErr)
					continue
				}
				if err != nil {
					return err
				}
			} else {
				if h.verbose {
					fmt.Printf("文件 %s 已存在，跳过\n", info.Name())
				}
			}
		}
		if !info.IsDir() {
			h.uploadProcess.EndFile(p, info)
		}
	}

	if filePath != h.srcPath {
		h.uploadProcess.StartDir(filePath, uint64(len(infos)))
	}
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	_, err = ReqRespWithConn(h.conn, rpcutil.CommandUploadEndDir, &pb.UploadEndDirReq{
		DriverId: h.driverId,
		DirPath:  dirPath,
	}, nil)
	if err != nil {
		return err
	}
	if filePath != h.srcPath {
		h.uploadProcess.EndDir(filePath, dirInfo)
	}

	return nil
}

func (h *uploadHandlersV3) copyFile(conn net.Conn, f *os.File, name string, size int64) error {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	if h.verbose {
		fmt.Printf("开始拷贝文件内容 %s，大小为 %s\n", name, humanize.IBytes(uint64(size)))
	}
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
	if h.verbose {
		fmt.Printf("拷贝文件内容完成 %s，大小为 %s\n", name, humanize.IBytes(uint64(n)))
	}
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

func (h *uploadHandlersV3) uploadFile(ctx context.Context, filePath string, info os.FileInfo) (hash string, fileErr error, err error) {
	f, fileErr := os.Open(filePath)
	if fileErr != nil {
		return
	}
	defer func() {
		f.Close()
	}()
	hash, fileErr = h.getHash(f)
	if fileErr != nil {
		return
	}

	select {
	case <-ctx.Done():
		err = context.Canceled
		return
	default:
	}
	if h.verbose {
		fmt.Printf("文件 %s 的哈希值为 %s\n", filePath, hash)
	}
	dirPath, err := h.formatPath(filepath.Dir(filePath))
	if err != nil {
		h.uploadProcess.OnFileError(filePath, err)
		return
	}

	modifyTime := uint64(info.ModTime().UnixNano())
	status, err := ReqRespWithConn(h.conn, rpcutil.CommandUploadV3File, &pb.UploadFileV3{
		DriverId:       h.driverId,
		DirPath:        dirPath,
		Name:           info.Name(),
		Hash:           hash,
		Mode:           uint64(info.Mode()),
		Size:           uint64(info.Size()),
		CreateTime:     modifyTime,
		ModifyTime:     modifyTime,
		ChangeTime:     modifyTime,
		AccessTime:     modifyTime,
		UploadDeviceId: h.uploadDeviceId,
		UploadTime:     h.uploadTime,
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
		//println("encoder", len(h.encoder), h.encoder)
		err = rpcutil.WriteString(h.conn, h.encoder)
		if err != nil {
			return
		}

		err = h.copyFile(h.conn, f, info.Name(), info.Size())
		if err != nil {
			return
		}

		_, _, err = rpcutil.ReadStatus(h.conn)
		if err != nil {
			return
		}
		return
	} else {
		if h.verbose {
			fmt.Printf("文件 %s 已存在，跳过\n", info.Name())
		}
	}
	return
}
