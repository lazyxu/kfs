package client

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"net"
	"os"

	"github.com/lazyxu/kfs/pb"
)

type uploadHandlers struct {
	core.DefaultWalkHandlers[core.FileResp]
	uploadProcess    core.UploadProcess
	concurrent       int
	encoder          string
	verbose          bool
	socketServerAddr string
	conns            []net.Conn
	files            []*os.File
}

func (h *uploadHandlers) StartWorker(ctx context.Context, index int) {
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

func (h *uploadHandlers) reconnect(ctx context.Context, index int) {
	h.conns[index].Close()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.socketServerAddr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	h.conns[index] = conn
}

func (h *uploadHandlers) EndWorker(ctx context.Context, index int) {
	h.conns[index].Close()
}

func (h *uploadHandlers) OnFileError(filePath string, info os.FileInfo, err error) {
	h.uploadProcess.OnFileError(-1, filePath, info, err)
}

func (h *uploadHandlers) StackSizeHandler(size int) {
	h.uploadProcess.StackSizeHandler(size)
}

func (h *uploadHandlers) StartFile(ctx context.Context, index int, filePath string, info os.FileInfo) {
	h.uploadProcess.StartFile(index, filePath, info)
}

func (h *uploadHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []core.FileResp) (fileResp core.FileResp) {
	var err error
	defer func() {
		if err != nil {
			h.uploadProcess.OnFileError(index, filePath, info, err)
		}
	}()
	fileResp.Info = info
	if info.Mode().IsRegular() {
		file, info, err, notExist := h.uploadFile(ctx, index, filePath)
		if err != nil {
			return
		}
		h.uploadProcess.EndFile(index, filePath, info, !notExist)
		fileResp.FileOrDir = file
		return
	} else if info.IsDir() {
		dirItems := make([]*pb.DirItem, len(children))
		for i, child := range children {
			if child.FileOrDir == nil {
				continue
			}
			modifyTime := uint64(child.Info.ModTime().UnixNano())
			dirItems[i] = &pb.DirItem{
				Hash:       child.FileOrDir.Hash(),
				Name:       child.Info.Name(),
				Mode:       uint64(child.Info.Mode()),
				Size:       child.FileOrDir.Size(),
				Count:      child.FileOrDir.Count(),
				TotalCount: child.FileOrDir.TotalCount(),
				CreateTime: modifyTime,
				ModifyTime: modifyTime,
				ChangeTime: modifyTime,
				AccessTime: modifyTime,
			}
		}

		var resp pb.UploadResp
		err = ReqResp(h.socketServerAddr, rpcutil.CommandUploadDirItem, &pb.UploadReq{
			Dir: &pb.UploadReqDir{DirItem: dirItems},
		}, &resp)
		if err != nil {
			return
		}
		h.uploadProcess.EndFile(index, filePath, info, resp.Exist)
		fileResp.FileOrDir = dao.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return
	}
	return
}

func (h *uploadHandlers) EnqueueFile(info os.FileInfo) {
	h.uploadProcess.EnqueueFile(info)
}
