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
}

func (h *uploadHandlers) StartWorker(ctx context.Context, index int) {
	conn, err := net.Dial("tcp", h.socketServerAddr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	h.conns[index] = conn
}

func (h *uploadHandlers) EndWorker(ctx context.Context, index int) {
	h.conns[index].Close()
}

func (h *uploadHandlers) ErrHandler(filePath string, err error) {
	h.uploadProcess.ErrHandler(filePath, err)
}

func (h *uploadHandlers) StackSizeHandler(size int) {
	h.uploadProcess.StackSizeHandler(size)
}

func (h *uploadHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []core.FileResp) (fileResp core.FileResp) {
	var err error
	defer func() {
		if err != nil {
			h.ErrHandler(filePath, err)
		}
	}()
	fileResp.Info = info
	if info.Mode().IsRegular() {
		file, err := h.uploadFile(ctx, index, filePath)
		if err != nil {
			return
		}
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
		fileResp.FileOrDir = dao.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return
	}
	return
}
