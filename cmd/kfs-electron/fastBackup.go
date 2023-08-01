package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/lazyxu/kfs/rpc/client"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/core"
)

type FastBackupWalker struct {
	FileSizeResp
	core.DefaultWalkHandlers[CountAndSize]
	req    WsReq
	onResp func(finished bool, data interface{}) error
	tick   <-chan time.Time
}

func (w *FastBackupWalker) StackSizeHandler(size int) {
	w.StackSize = size
}

func (w *FastBackupWalker) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []CountAndSize) CountAndSize {
	var count int64 = 1
	var size int64
	if info.IsDir() {
		atomic.AddInt64(&w.DirCount, 1)
		for _, child := range children {
			count += child.Count
			size += child.Size
		}
	} else {
		count = 1
		size = info.Size()
		atomic.AddInt64(&w.FileCount, 1)
		atomic.AddInt64(&w.FileSize, info.Size())
	}

	select {
	case <-w.tick:
		fmt.Printf("tick: %+v\n", w.FileSizeResp)
		err := w.onResp(false, w.FileSizeResp)
		if err != nil {
			fmt.Printf("%+v %+v\n", w.req, err)
		}
	case <-ctx.Done():
	default:
	}
	return CountAndSize{
		Count: count,
		Size:  size,
	}
}

func (p *WsProcessor) fastBackup(ctx context.Context, req WsReq, srcPath string, serverAddr string, branchName string, dstPath string) error {
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
	concurrent := 1
	encoder := ""

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}
	commit, branch, err := fs.Upload(ctx, branchName, dstPath, srcPath, core.UploadConfig{
		Encoder:    encoder,
		Concurrent: concurrent,
		Verbose:    true,
	})
	if err != nil {
		return p.err(req, err)
	}
	fmt.Printf("hash=%s, commitId=%d, size=%s, count=%d\n", commit.Hash[:4], branch.CommitId, humanize.Bytes(branch.Size), branch.Count)
	return p.ok(req, true, branch)
}
