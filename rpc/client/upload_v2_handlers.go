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

	"github.com/lazyxu/kfs/pb"
)

type uploadHandlersV2 struct {
	core.DefaultWalkByLevelHandlers
	uploadProcess    core.UploadProcess
	concurrent       int
	encoder          string
	verbose          bool
	socketServerAddr string
	conns            []net.Conn
	files            []*os.File
	driverName       string
	srcPath          string
	dstPath          string
}

func (h *uploadHandlersV2) StartWorker(ctx context.Context, index int) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.socketServerAddr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	h.conns[index] = conn
	go func() {
		<-ctx.Done()
		if h.files[index] != nil {
			h.files[index].Close()
		}
	}()
}

func (h *uploadHandlersV2) reconnect(ctx context.Context, index int) {
	h.conns[index].Close()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.socketServerAddr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	h.conns[index] = conn
}

func (h *uploadHandlersV2) EndWorker(ctx context.Context, index int) {
	h.conns[index].Close()
}

func (h *uploadHandlersV2) OnFileError(filePath string, index int, info os.FileInfo, err error) {
	h.uploadProcess.OnFileError(index, filePath, info, err)
}

func (h *uploadHandlersV2) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo) (err error) {
	h.uploadProcess.StartFile(index, filePath, info)

	var relPath string
	relPath, err = h.formatPath(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			h.reconnect(ctx, index)
			return
		}
	}()
	conn := h.conns[index]

	if info.Mode().IsRegular() {
		var exist bool
		exist, err = h.uploadFile(ctx, conn, index, filePath, info, relPath)
		if err != nil {
			return err
		}
		h.uploadProcess.EndFile(index, filePath, info, exist)
	} else if info.IsDir() {
		modifyTime := uint64(info.ModTime().UnixNano())
		_, err = ReqRespWithConn(conn, rpcutil.CommandUploadV2Dir, &pb.UploadReqV2{
			DriverName: h.driverName,
			DstPath:    h.dstPath,
			RelPath:    relPath,
			Hash:       "",
			Mode:       uint64(info.Mode()),
			Size:       0,
			CreateTime: modifyTime,
			ModifyTime: modifyTime,
			ChangeTime: modifyTime,
			AccessTime: modifyTime,
		}, nil)
		if err != nil {
			return err
		}
		h.uploadProcess.EndFile(index, filePath, info, false)
		return nil
	}
	return nil
}

func (h *uploadHandlersV2) AddToWorkList(info os.FileInfo) {
	h.uploadProcess.PushFile(info)
}

func (h *uploadHandlersV2) HasEnqueuedAll() {
	h.uploadProcess.HasPushedAllToStack()
}

func (h *uploadHandlersV2) copyFile(conn net.Conn, f *os.File, size int64) error {
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

func (h *uploadHandlersV2) getHash(f *os.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (h *uploadHandlersV2) formatPath(filePath string) (string, error) {
	rel, err := filepath.Rel(h.srcPath, filePath)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

func (h *uploadHandlersV2) uploadFile(ctx context.Context, conn net.Conn, index int, filePath string, info os.FileInfo, relPath string) (exist bool, err error) {
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

	println(conn.RemoteAddr().String(), "hash", len(hash), hash)

	modifyTime := uint64(info.ModTime().UnixNano())
	status, err := ReqRespWithConn(conn, rpcutil.CommandUploadV2File, &pb.UploadReqV2{
		DriverName: h.driverName,
		DstPath:    h.dstPath,
		RelPath:    relPath,
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
