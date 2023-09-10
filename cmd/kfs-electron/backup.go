package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func apiBackupTask(c echo.Context) error {
	list, err := db.ListBackupTask(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, list)
}

type Client struct {
	sseChannel     chan string
	sseJsonChannel chan interface{}
}

var clients sync.Map // map[*http.Request]*Client

func apiEventBackupTask(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	sseChannel := make(chan string)
	sseJsonChannel := make(chan interface{})
	client := &Client{
		sseChannel:     sseChannel,
		sseJsonChannel: sseJsonChannel,
	}
	clients.Store(c, client)

	defer func() {
		close(sseChannel)
		close(sseJsonChannel)
		clients.Delete(c)
		fmt.Println("Closing connection")
	}()

	obj, err := getTaskInfos(c.Request().Context())
	if err != nil {
		return err
	}
	data, err := json.Marshal(obj)
	if err != nil {
		log.Panicf("invalid obj: %+v\n", obj)
	}
	fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
	c.Response().Flush()

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

type TaskInfos struct {
	List         []BackupTask                  `json:"list"`
	RunningTasks map[string]*RunningBackupTask `json:"runningTasks"`
}

func getTaskInfos(ctc context.Context) (info TaskInfos, err error) {
	list, err := db.ListBackupTask(ctc)
	if err != nil {
		return
	}
	runningTasksMutex.RLock()
	defer runningTasksMutex.RUnlock()

	return TaskInfos{
		List:         list,
		RunningTasks: runningTasks,
	}, nil
}

func noteTaskListToClients() {
	clients.Range(func(key, value any) bool {
		c := key.(echo.Context)
		client := value.(*Client)
		obj, err := getTaskInfos(c.Request().Context())
		if err != nil {
			c.Logger().Error(err)
			return true
		}
		client.sseJsonChannel <- obj
		return true
	})
}

func apiNewBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	srcPath := c.QueryParam("srcPath")
	driverName := c.QueryParam("driverName")
	dstPath := c.QueryParam("dstPath")
	encoder := c.QueryParam("encoder")
	concurrentStr := c.QueryParam("concurrent")
	concurrent, err := strconv.Atoi(concurrentStr)
	if err != nil {
		return err
	}
	err = upsertBackup(c.Request().Context(), db, name, description, srcPath, driverName, dstPath, encoder, concurrent)
	if err != nil {
		return err
	}
	noteTaskListToClients()
	return c.String(http.StatusOK, "")
}

func apiDeleteBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	err := db.DeleteBackupTask(c.Request().Context(), name)
	if err != nil {
		return err
	}
	noteTaskListToClients()
	return c.String(http.StatusOK, "")
}

type RunningBackupTask struct {
	cancel       context.CancelFunc
	Status       int   `json:"status"`
	LastDoneTime int64 `json:"lastDoneTime"`
}

var (
	StatusIdle        = 0
	StatusWaitRunning = 1
	StatusRunning     = 2
	StatusFinished    = 3
	StatusCanceled    = 4
	StatusError       = 5
)

var runningTasks = make(map[string]*RunningBackupTask)
var runningTasksMutex = &sync.RWMutex{}

func apiStartBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	startStr := c.QueryParam("start")
	serverAddr := c.QueryParam("serverAddr")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	task, err := db.GetBackupTask(c.Request().Context(), name)
	if err != nil {
		return err
	}
	runningTasksMutex.Lock()
	defer runningTasksMutex.Unlock()
	runningTask, exist := runningTasks[task.Name]
	if !exist {
		runningTask = &RunningBackupTask{
			cancel: nil,
			Status: StatusIdle,
		}
		runningTasks[task.Name] = runningTask
		if !start {
			return c.String(http.StatusOK, "")
		}
		tryStartBackup(task, runningTask, serverAddr)
	} else {
		if !start {
			runningTask.cancel()
			return c.String(http.StatusOK, "")
		}
		tryStartBackup(task, runningTask, serverAddr)
	}
	return c.String(http.StatusOK, "")
}

func tryStartBackup(task BackupTask, runningTask *RunningBackupTask, serverAddr string) {
	if runningTask.Status == StatusWaitRunning || runningTask.Status == StatusRunning {
		return
	}
	runningTask.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	runningTask.cancel = cancel
	go func() {
		setTaskStatus(task.Name, StatusRunning)
		err := eventSourceBackup(ctx, task.Name, task.Description, task.SrcPath, serverAddr, task.DriverName, task.DstPath, task.Encoder, task.Concurrent)
		if err == nil {
			setTaskStatus(task.Name, StatusFinished)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			setTaskStatus(task.Name, StatusCanceled)
			return
		}
		setTaskStatus(task.Name, StatusError)
	}()
}

func setTaskStatus(name string, status int) {
	runningTasksMutex.Lock()
	runningTask := runningTasks[name]
	runningTask.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		// TODO: save it to db.
		runningTask.LastDoneTime = time.Now().UnixNano()
	}
	runningTasksMutex.Unlock()
	noteTaskListToClients()
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
	select {
	case <-ctx.Done():
		fmt.Println("backup canceled")
		return context.DeadlineExceeded
	case <-time.After(time.Second * 10):
	}
	fmt.Println("backup finish")
	return nil
	//fs := &client.RpcFs{
	//	SocketServerAddr: serverAddr,
	//}

	//w := NewWebUploadProcess(ctx, req, concurrent, func(finished bool, data interface{}) error {
	//	return fmt.Printf(req, finished, data)
	//})
	//
	//err = fs.UploadV2(ctx, driverName, dstPath, srcPath, core.UploadConfig{
	//	UploadProcess: w,
	//	Encoder:       encoder,
	//	Concurrent:    concurrent,
	//	Verbose:       false,
	//})
	//if err != nil {
	//	return p.err(req, err)
	//}
	//for i := 0; i < concurrent; i++ {
	//	w.Done <- struct{}{}
	//}
	//for i := 0; i < concurrent; i++ {
	//	w.RespIfUpdated(i)
	//}
	//fmt.Printf("w=%+v\n", w)
	//return p.ok(req, true, WebBackupResp{
	//	Size: w.Size, FileCount: w.FileCount, DirCount: w.DirCount,
	//	TotalSize: w.TotalSize, TotalFileCount: w.TotalFileCount, TotalDirCount: w.TotalDirCount,
	//	Processes: w.Processes[:], PushedAllToStack: w.PushedAllToStack, Cost: time.Now().Sub(w.StartTime).Milliseconds(),
	//})
}
