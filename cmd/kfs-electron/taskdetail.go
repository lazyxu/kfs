package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	"log"
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

func apiEventBackupTaskDetail(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	name := c.Param("name")
	fmt.Println("New connection established", name)
	sseChannel := make(chan string)
	sseJsonChannel := make(chan interface{})
	addTaskDetailClient(name, c, &Client{
		sseChannel:     sseChannel,
		sseJsonChannel: sseJsonChannel,
	})

	defer func() {
		close(sseChannel)
		close(sseJsonChannel)
		deleteTaskDetailClient(name, c)
		fmt.Println("Closing connection", name)
	}()

	for {
		select {
		case msg := <-sseChannel:
			fmt.Fprintf(c.Response(), "data: %s\n\n", msg)
			c.Response().Flush()

		case obj := <-sseJsonChannel:
			data, err := json.Marshal(obj)
			if err != nil {
				log.Panicf("invalid obj: %+v\n", obj)
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
			c.Response().Flush()

		case <-c.Request().Context().Done():
			fmt.Println("Connection closed")
			return nil
		}
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

func eventSourceBackup(ctx context.Context, name, description, srcPath, serverAddr, driverName, dstPath, encoder string, concurrent int) error {
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

	err = fs.UploadV2(ctx, driverName, dstPath, srcPath, core.UploadConfig{
		UploadProcess: w,
		Encoder:       encoder,
		Concurrent:    concurrent,
		Verbose:       false,
	})
	if err != nil {
		return errTaskDetailToClients(name, err)
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

func NewWebUploadProcess(ctx context.Context, concurrent int, onResp func(finished bool, data interface{}) error) *WebUploadProcess {
	w := &WebUploadProcess{
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
