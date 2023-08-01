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
	core.DefaultWalkHandlers[FileResp]
	uploadProcess    core.UploadProcess
	concurrent       int
	encoder          string
	verbose          bool
	socketServerAddr string
	conns            []net.Conn
}

type FileResp struct {
	fileOrDir dao.IFileOrDir
	info      os.FileInfo
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

func (h *uploadHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []FileResp) (fileResp FileResp) {
	var err error
	defer func() {
		if err != nil {
			h.ErrHandler(filePath, err)
		}
	}()
	fileResp.info = info
	if info.Mode().IsRegular() {
		file, err := h.uploadFile(ctx, index, filePath)
		if err != nil {
			return
		}
		fileResp.fileOrDir = file
		return
	} else if info.IsDir() {
		dirItems := make([]*pb.DirItem, len(children))
		for i, child := range children {
			if child.fileOrDir == nil {
				continue
			}
			modifyTime := uint64(child.info.ModTime().UnixNano())
			dirItems[i] = &pb.DirItem{
				Hash:       child.fileOrDir.Hash(),
				Name:       child.info.Name(),
				Mode:       uint64(child.info.Mode()),
				Size:       child.fileOrDir.Size(),
				Count:      child.fileOrDir.Count(),
				TotalCount: child.fileOrDir.TotalCount(),
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
		fileResp.fileOrDir = dao.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return
	}
	return
}
