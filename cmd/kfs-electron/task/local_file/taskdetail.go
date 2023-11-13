package local_file

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-electron/backup"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var taskDetailClients = make(map[string]map[echo.Context]*Client)
var taskDetailClientsMutex = &sync.RWMutex{}

func addTaskDetailClient(name string, c echo.Context, client *Client) {
	taskDetailClientsMutex.Lock()
	defer taskDetailClientsMutex.Unlock()
	if clients, ok := taskDetailClients[name]; ok {
		clients[c] = client
	} else {
		taskDetailClients[name] = make(map[echo.Context]*Client)
		taskDetailClients[name][c] = client
	}
	taskDetailClients[name][c] = client
}

func deleteTaskDetailClient(name string, c echo.Context) {
	taskDetailClientsMutex.Lock()
	defer taskDetailClientsMutex.Unlock()
	if clients, ok := taskDetailClients[name]; ok {
		delete(clients, c)
	}
}

func noteTaskDetailToClients(name string, resp WsResp) {
	taskDetailClientsMutex.Lock()
	defer taskDetailClientsMutex.Unlock()
	if clients, ok := taskDetailClients[name]; ok {
		for _, client := range clients {
			client.sseJsonChannel <- resp
		}
	}
}

func okTaskDetailToClients(name string, finished bool, data interface{}) error {
	var resp WsResp
	resp.Finished = finished
	resp.Data = data
	noteTaskDetailToClients(name, resp)
	return nil
}

func errTaskDetailToClients(name string, err error) error {
	var resp WsResp
	resp.Finished = true
	resp.ErrMsg = err.Error()
	noteTaskDetailToClients(name, resp)
	return nil
}

const (
	LevelInfo  = 0
	LevelWarn  = 1
	LevelError = 2
)

type BackupRecord struct {
	Time    int64  `json:"time"`
	Level   int    `json:"level"`
	Content string `json:"content"`
}

type BackupLogs struct {
	Finished bool           `json:"finished"`
	Records  []BackupRecord `json:"records"`

	Size      uint64 `json:"size"`
	FileCount uint64 `json:"fileCount"`
	DirCount  uint64 `json:"dirCount"`

	TotalSize      uint64 `json:"totalSize"`
	TotalFileCount uint64 `json:"totalFileCount"`
	TotalDirCount  uint64 `json:"totalDirCount"`

	Cost int64 `json:"cost"`

	ErrMsg string `json:"errMsg,omitempty"`
}

func (d *DriverLocalFile) eventSourceBackup(ctx context.Context, name, description, srcPath, serverAddr string, driverId uint64, dstPath, encoder string, concurrent int) error {
	if !filepath.IsAbs(srcPath) {
		return errors.New("请输入绝对路径")
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("源目录不存在")
	}
	fmt.Println("backup start")

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}

	w := NewWebUploadProcess(ctx, concurrent, func(finished bool, data interface{}) error {
		return okTaskDetailToClients(name, finished, data)
	})

	err = fs.UploadV2(ctx, driverId, dstPath, srcPath, core.UploadConfig{
		UploadProcess: w,
		Encoder:       encoder,
		Concurrent:    concurrent,
		Verbose:       false,
	})
	if err != nil {
		return err
	}
	for i := 0; i < concurrent; i++ {
		w.Done <- struct{}{}
	}
	for i := 0; i < concurrent; i++ {
		w.RespIfUpdated(i)
	}
	fmt.Printf("w=%+v\n", w)
	fmt.Println("backup finish")
	return okTaskDetailToClients(name, true, WebBackupResp{
		Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
		TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
		Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack, Cost: time.Now().Sub(w.StartTime).Milliseconds(),
	})
}

func NewWebUploadProcess(ctx context.Context, concurrent int, onResp func(finished bool, data interface{}) error) *backup.WebUploadProcess {
	w := &backup.WebUploadProcess{
		ctx:       ctx,
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
					return
				case <-ctx.Done():
					<-w.Done
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
