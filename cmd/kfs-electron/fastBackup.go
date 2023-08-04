package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/dao"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/lazyxu/kfs/rpc/client"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/core"
)

type WebUploadProcess struct {
	Size           uint64
	FileCount      uint64
	DirCount       uint64
	TotalSize      uint64
	TotalFileCount uint64
	TotalDirCount  uint64
	ctx            context.Context
	req            WsReq
	onResp         func(finished bool, data interface{}) error
	tick           <-chan time.Time
}

type WebBackupResp struct {
	FilePath       string     `json:"filePath"`
	Size           uint64     `json:"size"`
	FileCount      uint64     `json:"fileCount"`
	DirCount       uint64     `json:"dirCount"`
	TotalSize      uint64     `json:"totalSize"`
	TotalFileCount uint64     `json:"totalFileCount"`
	TotalDirCount  uint64     `json:"totalDirCount"`
	Err            error      `json:"err"`
	Exist          bool       `json:"exist"`
	Branch         dao.Branch `json:"branch"`
}

func (w *WebUploadProcess) Show(p *core.Process) {
}

func (w *WebUploadProcess) StackSizeHandler(size int) {
	w.Show(&core.Process{
		StackSize: size,
	})
}

func (w *WebUploadProcess) New(srcPath string, concurrent int, conns []net.Conn) core.UploadProcess {
	return w
}

func (w *WebUploadProcess) Close(resp core.FileResp, err error) {
}

func (w *WebUploadProcess) OnFileError(filePath string, info os.FileInfo, err error) {
	println(filePath+":", err.Error())
	e := w.onResp(false, WebBackupResp{
		FilePath: filePath, Err: err,
		Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
		TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
	})
	if e != nil {
		fmt.Printf("%+v %+v\n", w.req, e)
	}
}

func (w *WebUploadProcess) EndFile(filePath string, info os.FileInfo, exist bool) {
	if info.IsDir() {
		w.DirCount++
	} else {
		w.FileCount++
		w.Size += uint64(info.Size())
	}
	select {
	case <-w.ctx.Done():
		return
	default:
	}
	select {
	case <-w.tick:
		e := w.onResp(false, WebBackupResp{
			FilePath: filePath, Exist: exist,
			Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
			TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
		})
		if e != nil {
			fmt.Printf("%+v %+v\n", w.req, e)
		}
	default:
	}
}

func (w *WebUploadProcess) EnqueueFile(info os.FileInfo) {
	if info.IsDir() {
		w.TotalDirCount++
	} else {
		w.TotalFileCount++
		w.TotalSize += uint64(info.Size())
	}
}

func (w *WebUploadProcess) Verbose() bool {
	return true
}

func (p *WsProcessor) fastBackup(ctx context.Context, req WsReq, srcPath string, serverAddr string, branchName string, dstPath string, concurrent int, encoder string, verbose bool) error {
	if !filepath.IsAbs(srcPath) {
		return p.err(req, errors.New("请输入绝对路径"))
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		return p.err(req, err)
	}
	if !info.IsDir() {
		return p.err(req, errors.New("请输入一个目录"))
	}

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}

	var uploadProcess core.UploadProcess
	if verbose {
		uploadProcess = &WebUploadProcess{
			ctx:  ctx,
			req:  req,
			tick: time.Tick(time.Millisecond * 500),
			onResp: func(finished bool, data interface{}) error {
				return p.ok(req, finished, data)
			},
		}
	} else {
		uploadProcess = &core.EmptyUploadProcess{}
	}

	commit, branch, err := fs.Upload(ctx, branchName, dstPath, srcPath, core.UploadConfig{
		UploadProcess: uploadProcess,
		Encoder:       encoder,
		Concurrent:    concurrent,
		Verbose:       false,
	})
	if err != nil {
		return p.err(req, err)
	}
	fmt.Printf("hash=%s, commitId=%d, size=%s, count=%d\n", commit.Hash[:4], branch.CommitId, humanize.Bytes(branch.Size), branch.Count)
	return p.ok(req, true, WebBackupResp{Branch: branch})
}
