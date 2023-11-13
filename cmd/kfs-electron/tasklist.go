package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/cmd/kfs-electron/backup"
	"github.com/lazyxu/kfs/cmd/kfs-electron/db/gosqlite"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
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

var taskListClients sync.Map // map[*http.Request]*Client

func apiEventBackupTask(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")
	sseChannel := make(chan string)
	sseJsonChannel := make(chan interface{})
	taskListClients.Store(c, &Client{
		sseChannel:     sseChannel,
		sseJsonChannel: sseJsonChannel,
	})

	defer func() {
		close(sseChannel)
		close(sseJsonChannel)
		taskListClients.Delete(c)
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
	List         []gosqlite.BackupTask         `json:"list"`
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
	taskListClients.Range(func(key, value any) bool {
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
	err = backup.UpsertBackup(c.Request().Context(), db, name, description, srcPath, driverName, dstPath, encoder, concurrent)
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
	StatusIdle         = 0
	StatusWaitRunning  = 1
	StatusRunning      = 2
	StatusFinished     = 3
	StatusCanceled     = 4
	StatusError        = 5
	StatusWaitCanceled = 6
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
			setTaskStatus(task.Name, StatusWaitCanceled)
			runningTask.cancel()
			return c.String(http.StatusOK, "")
		}
		tryStartBackup(task, runningTask, serverAddr)
	}
	return c.String(http.StatusOK, "")
}

func tryStartBackup(task gosqlite.BackupTask, runningTask *RunningBackupTask, serverAddr string) {
	if runningTask.Status == StatusWaitRunning || runningTask.Status == StatusRunning {
		return
	}
	runningTask.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	runningTask.cancel = cancel
	go func() {
		setTaskStatus(task.Name, StatusRunning)
		err := eventSourceBackup(ctx, task.Name, task.Description, task.SrcPath, serverAddr, task.DriverId, task.DstPath, task.Encoder, task.Concurrent)
		if err == nil {
			setTaskStatus(task.Name, StatusFinished)
			return
		}
		if errors.Is(err, context.Canceled) {
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
		runningTask.cancel = nil
	}
	runningTasksMutex.Unlock()
	noteTaskListToClients()
}
