package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type FastScanWalker struct {
	FileSizeResp
	core.DefaultWalkHandlers[CountAndSize]
	req    WsReq
	onResp func(finished bool, data interface{}) error
	tick   <-chan time.Time
}

func (w *FastScanWalker) StackSizeHandler(size int) {
	w.StackSize = size
}

func (w *FastScanWalker) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []CountAndSize) CountAndSize {
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

func (p *WsProcessor) fastScan(ctx context.Context, req WsReq, backupDir string) error {
	if !filepath.IsAbs(backupDir) {
		return p.err(req, errors.New("请输入绝对路径"))
	}
	info, err := os.Lstat(backupDir)
	if err != nil {
		return p.err(req, err)
	}
	if !info.IsDir() {
		return p.err(req, errors.New("请输入一个目录"))
	}
	w := FastScanWalker{
		req:  req,
		tick: time.Tick(time.Millisecond * 500),
		onResp: func(finished bool, data interface{}) error {
			return p.ok(req, finished, data)
		},
	}
	err = p.ok(req, false, w.FileSizeResp)
	if err != nil {
		return err
	}
	_, err = core.Walk[CountAndSize](ctx, backupDir, 15, &w)
	if err != nil {
		return p.err(req, err)
	}
	err = p.ok(req, false, w.FileSizeResp)
	if err != nil {
		return err
	}
	return p.ok(req, true, w.FileSizeResp)
}
