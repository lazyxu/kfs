package metadata

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
	"github.com/lazyxu/kfs/rpc/server"
)

type Client struct {
	ch chan TaskInfo
}

func (c *Client) Chan() chan TaskInfo {
	return c.ch
}

func (c *Client) Message() TaskInfo {
	return taskInfo
}

var s = &core.EventServer[TaskInfo]{
	Clients: make(map[echo.Context]core.EventClient[TaskInfo]),
	NewClient: func(c echo.Context, kfsCore *core.KFS) (core.EventClient[TaskInfo], error) {
		return &Client{
			ch: make(chan TaskInfo),
		}, nil
	},
}

func ApiEvent(c echo.Context, kfsCore *core.KFS) error {
	return s.Handle(c, kfsCore)
}

var (
	StatusIdle           = 0
	StatusFinished       = 1
	StatusCanceled       = 2
	StatusError          = 3
	StatusWaitRunning    = 4
	StatusWaitCanceled   = 5
	StatusRunningCollect = 6
	StatusRunningAnalyze = 7
)

type TaskInfo struct {
	cancel       context.CancelFunc
	Status       int      `json:"status"`
	LastDoneTime int64    `json:"lastDoneTime"`
	Cnt          int      `json:"cnt"`
	Total        int      `json:"total"`
	Errors       []string `json:"errors"`
}

var taskInfo = TaskInfo{
	cancel: nil,
	Status: StatusIdle,
	Cnt:    0,
	Total:  0,
	Errors: make([]string, 0),
}

var mutex = &sync.RWMutex{}

func setTaskStatus(status int) {
	mutex.Lock()
	taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		taskInfo.cancel = nil
		taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunningCollect {
		taskInfo.Errors = make([]string, 0)
		taskInfo.Cnt = 0
		taskInfo.Total = 0
	}
	mutex.Unlock()
	s.SendAll()
}

func setTaskStatusWithLock(status int) {
	taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		taskInfo.cancel = nil
		taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunningCollect {
		taskInfo.Errors = make([]string, 0)
		taskInfo.Cnt = 0
		taskInfo.Total = 0
	}
	s.SendAll()
}

func setTaskTotal(total int) {
	mutex.Lock()
	taskInfo.Status = StatusRunningAnalyze
	taskInfo.Errors = make([]string, 0)
	taskInfo.Cnt = 0
	taskInfo.Total = total
	mutex.Unlock()
	s.SendAll()
}

func addTaskCnt() {
	mutex.Lock()
	taskInfo.Cnt++
	mutex.Unlock()
	s.SendAll()
}

func addTaskError(err error) {
	mutex.Lock()
	taskInfo.Errors = append(taskInfo.Errors, err.Error())
	mutex.Unlock()
	s.SendAll()
}

func StartOrStop(kfsCore *core.KFS, start bool, all bool) {
	mutex.Lock()
	defer mutex.Unlock()
	if !start {
		setTaskStatusWithLock(StatusWaitCanceled)
		taskInfo.cancel()
		return
	}
	if taskInfo.Status == StatusWaitRunning ||
		taskInfo.Status == StatusRunningCollect ||
		taskInfo.Status == StatusRunningAnalyze ||
		taskInfo.Status == StatusWaitCanceled {
		return
	}
	taskInfo.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	taskInfo.cancel = cancel
	go func() {
		err := analyze(ctx, kfsCore, all)
		if err == nil {
			setTaskStatus(StatusFinished)
			return
		}
		if errors.Is(err, context.Canceled) {
			setTaskStatus(StatusCanceled)
			return
		}
		setTaskStatus(StatusError)
	}()
}

func analyze(ctx context.Context, kfsCore *core.KFS, all bool) (err error) {
	setTaskStatus(StatusRunningCollect)
	var fileList []sqlc.File
	if all {
		fileList, err = kfsCore.Q.ListFile(ctx)
	} else {
		fileList, err = kfsCore.Q.ListFileExpectFileType(ctx)
	}
	if err != nil {
		return err
	}
	setTaskTotal(len(fileList))
	for _, f := range fileList {
		if all {
			err = server.Analyze(ctx, kfsCore, f.ID, f.Size, f.Md5, f.Sha256)
		} else {
			err = server.AnalyzeIfNoFileType(ctx, kfsCore, f.ID, f.Size, f.Md5, f.Sha256)
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			addTaskError(err)
			continue
		}
		addTaskCnt()
	}
	if err != nil {
		return err
	}
	return nil
}
