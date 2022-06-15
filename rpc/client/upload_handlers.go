package client

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type uploadHandlers struct {
	core.DefaultWalkHandlers[fileResp]
	c             pb.KoalaFSClient
	uploadProcess core.UploadProcess
	concurrent    int
	encoder       string
	verbose       bool
	ch            chan *Process
	conns         []net.Conn
}

type fileResp struct {
	fileOrDir sqlite.FileOrDir
	info      os.FileInfo
}

func (h *uploadHandlers) BeforeFileHandler(ctx context.Context, index int) {
	conn, err := net.Dial("tcp", "127.0.0.1:1124")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	h.conns[index] = conn
}

func (h *uploadHandlers) AfterFileHandler(ctx context.Context, index int) {
	h.conns[index].Close()
}

func (h *uploadHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []fileResp) (fileResp fileResp) {
	var err error
	defer func() {
		if err != nil {
			h.ErrHandler(filePath, err)
		}
	}()
	fileResp.info = info
	if info.Mode().IsRegular() {
		h.uploadProcess = h.uploadProcess.New(int(info.Size()), filepath.Base(filePath))
		defer h.uploadProcess.Close()
		fileResp.fileOrDir, err = h.uploadFile(ctx, index, filePath)
		if err != nil {
			return
		}
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
		var client pb.KoalaFS_UploadClient
		client, err = h.c.Upload(ctx)
		if err != nil {
			return
		}
		err = client.Send(&pb.UploadReq{
			Dir: &pb.UploadReqDir{DirItem: dirItems},
		})
		var resp *pb.UploadResp
		resp, err = client.Recv()
		if err != nil {
			return
		}
		err = client.CloseSend()
		if err != nil {
			return
		}
		fileResp.fileOrDir = sqlite.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return
	}
	return
}
