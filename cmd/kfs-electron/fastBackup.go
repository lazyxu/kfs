package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/dao"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/lazyxu/kfs/rpc/client"

	"github.com/lazyxu/kfs/core"
)

type WebUploadProcess struct {
	Size             uint64
	FileCount        uint64
	DirCount         uint64
	TotalSize        uint64
	TotalFileCount   uint64
	TotalDirCount    uint64
	Processes        []Process
	PushedAllToStack bool

	StartTime time.Time
	ctx       context.Context
	Done      chan struct{}
	req       WsReq
	onResp    func(finished bool, data interface{}) error
}

type Process struct {
	updated  atomic.Bool
	FilePath string `json:"filePath"`
	Size     uint64 `json:"size"`
	Status   int    `json:"status"`
}

type WebBackupResp struct {
	Size             uint64    `json:"size"`
	FileCount        uint64    `json:"fileCount"`
	DirCount         uint64    `json:"dirCount"`
	TotalSize        uint64    `json:"totalSize"`
	TotalFileCount   uint64    `json:"totalFileCount"`
	TotalDirCount    uint64    `json:"totalDirCount"`
	Processes        []Process `json:"processes"`
	PushedAllToStack bool      `json:"pushedAllToStack"`
	Cost             int64     `json:"cost"`

	FilePath string     `json:"filePath"`
	ErrMsg   string     `json:"errMsg"`
	Branch   dao.Branch `json:"branch"`
}

const (
	StatusUploading = iota
	StatusExist
	StatusUploaded
)

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

func (w *WebUploadProcess) StartFile(index int, filePath string, info os.FileInfo) {
	w.Processes[index] = Process{FilePath: filePath, Size: uint64(info.Size()), Status: StatusUploading}
	w.Processes[index].updated.Store(true)
}

func (w *WebUploadProcess) OnFileError(index int, filePath string, info os.FileInfo, err error) {
	if index != -1 {
		w.Processes[index] = Process{}
	}
	println(filePath+":", err.Error())
	e := w.onResp(false, WebBackupResp{
		FilePath: filePath, ErrMsg: err.Error(),
		Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
		TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
		Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack,
	})
	if e != nil {
		fmt.Printf("%+v %+v\n", w.req, e)
	}
}

func (w *WebUploadProcess) EndFile(index int, filePath string, info os.FileInfo, exist bool) {
	if info.IsDir() {
		w.DirCount++
	} else {
		w.FileCount++
		w.Size += uint64(info.Size())
	}
	if w.Processes[index].FilePath != filePath {
		panic("w.Processes[index].FilePath != filePath")
	}
	w.Processes[index].Status = StatusUploaded
	if exist {
		w.Processes[index].Status = StatusExist
	}
	w.Processes[index].updated.Store(true)
}

func (w *WebUploadProcess) PushFile(info os.FileInfo) {
	if info.IsDir() {
		w.TotalDirCount++
	} else {
		w.TotalFileCount++
		w.TotalSize += uint64(info.Size())
	}
	w.Processes[0].updated.Store(true)
}

func (w *WebUploadProcess) HasPushedAllToStack() {
	w.PushedAllToStack = true
}

func (w *WebUploadProcess) Verbose() bool {
	return true
}

func (p *WsProcessor) fastBackup(ctx context.Context, req WsReq, srcPath string, serverAddr string, branchName string, dstPath string, concurrent int, encoder string) error {
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

	w := NewWebUploadProcess(ctx, req, concurrent, func(finished bool, data interface{}) error {
		return p.ok(req, finished, data)
	})

	err = fs.UploadV2(ctx, branchName, dstPath, srcPath, core.UploadConfig{
		UploadProcess: w,
		Encoder:       encoder,
		Concurrent:    concurrent,
		Verbose:       false,
	})
	if err != nil {
		return p.err(req, err)
	}
	w.Done <- struct{}{}
	for i := 0; i < concurrent; i++ {
		w.RespIfUpdated(i)
	}
	fmt.Printf("w=%+v\n", w)
	return p.ok(req, true, WebBackupResp{
		Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
		TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
		Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack, Cost: time.Now().Sub(w.StartTime).Milliseconds(),
	})
}

func NewWebUploadProcess(ctx context.Context, req WsReq, concurrent int, onResp func(finished bool, data interface{}) error) *WebUploadProcess {
	w := &WebUploadProcess{
		ctx:       ctx,
		req:       req,
		onResp:    onResp,
		Processes: make([]Process, concurrent),
		Done:      make(chan struct{}),
		StartTime: time.Now(),
	}
	for i := 0; i < concurrent; i++ {
		go func(i int) {
			for {
				select {
				case <-w.Done:
					w.Done <- struct{}{}
					w.Resp(i)
					return
				case <-ctx.Done():
					w.Resp(i)
					return
				default:
					w.Resp(i)
				}
				time.Sleep(time.Millisecond * 500)
			}
		}(i)
	}
	return w
}

func (w *WebUploadProcess) RespIfUpdated(i int) {
	if w.Processes[i].updated.CompareAndSwap(true, false) {
		e := w.onResp(false, WebBackupResp{
			Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
			TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
			Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack, Cost: time.Now().Sub(w.StartTime).Milliseconds(),
		})
		if e != nil {
			fmt.Printf("%+v %+v\n", w.req, e)
		}
	}
}

func (w *WebUploadProcess) Resp(i int) {
	w.Processes[i].updated.Store(false)
	e := w.onResp(false, WebBackupResp{
		Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
		TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
		Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack, Cost: time.Now().Sub(w.StartTime).Milliseconds(),
	})
	if e != nil {
		fmt.Printf("%+v %+v\n", w.req, e)
	}
}
