package client

import (
	"context"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/lazyxu/kfs/core"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Upload(ctx context.Context, branchName string, dstPath string, srcPath string, config core.UploadConfig) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return
	}
	defer conn.Close()

	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	handlers := &uploadHandlers{
		c:                c,
		uploadProcess:    config.UploadProcess,
		encoder:          config.Encoder,
		verbose:          config.Verbose,
		concurrent:       config.Concurrent,
		socketServerAddr: fs.SocketServerAddr,
		ch:               make(chan *Process),
		conns:            make([]net.Conn, config.Concurrent),
	}
	var wg sync.WaitGroup
	if config.Verbose {
		handlers.ch = make(chan *Process)
		wg.Add(1)
		go func() {
			handlers.handleProcess(srcPath)
			wg.Done()
		}()
	}
	walkResp, err := core.Walk[fileResp](ctx, srcPath, config.Concurrent, handlers)
	if config.Verbose {
		close(handlers.ch)
		wg.Wait()
	}
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
	return sqlite.Commit{
			Id:   resp.Branch.CommitId,
			Hash: resp.Branch.Hash,
		}, sqlite.Branch{
			Name:     branchName,
			CommitId: resp.Branch.CommitId,
			Size:     resp.Branch.Size,
			Count:    resp.Branch.Count,
		}, nil
}
