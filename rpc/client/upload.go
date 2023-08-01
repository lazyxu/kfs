package client

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"net"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Upload(ctx context.Context, branchName string, dstPath string, srcPath string, config core.UploadConfig) (commit dao.Commit, branch dao.Branch, err error) {

	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	handlers := &uploadHandlers{
		uploadProcess:    config.UploadProcess,
		encoder:          config.Encoder,
		verbose:          config.Verbose,
		concurrent:       config.Concurrent,
		socketServerAddr: fs.SocketServerAddr,
		conns:            make([]net.Conn, config.Concurrent),
	}
	handlers.uploadProcess = handlers.uploadProcess.New(srcPath, config.Concurrent, handlers.conns)
	walkResp, err := core.Walk[FileResp](ctx, srcPath, config.Concurrent, handlers)
	handlers.uploadProcess.Close()
	if err != nil {
		return
	}
	info, err := os.Stat(srcPath)
	if err != nil {
		return
	}
	fileOrDir := walkResp.fileOrDir
	modifyTime := uint64(info.ModTime().UnixNano())

	var resp pb.UploadResp
	err = ReqResp(fs.SocketServerAddr, rpcutil.CommandUploadDirItem, &pb.UploadReq{
		Root: &pb.UploadReqRoot{
			BranchName: branchName,
			Path:       dstPath,
			DirItem: &pb.DirItem{
				Hash:       fileOrDir.Hash(),
				Name:       filepath.Base(dstPath),
				Mode:       uint64(info.Mode()),
				Size:       fileOrDir.Size(),
				Count:      fileOrDir.Count(),
				TotalCount: fileOrDir.TotalCount(),
				CreateTime: modifyTime,
				ModifyTime: modifyTime,
				ChangeTime: modifyTime,
				AccessTime: modifyTime,
			},
		},
	}, &resp)
	if err != nil {
		return
	}
	return dao.Commit{
			Id:   resp.Branch.CommitId,
			Hash: resp.Branch.Hash,
		}, dao.Branch{
			Name:     branchName,
			CommitId: resp.Branch.CommitId,
			Size:     resp.Branch.Size,
			Count:    resp.Branch.Count,
		}, nil
}
