package metadata

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/common"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/server"
	"sync"
	"time"
)

var s = common.EventServer[TaskInfo]{
	Message: func() TaskInfo {
		return taskInfo
	},
}

func ApiEvent(c echo.Context) error {
	return s.Handle(c)
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

func StartOrStop(kfsCore *core.KFS, start bool) {
	mutex.Lock()
	defer mutex.Unlock()
	if !start {
		setTaskStatus(StatusWaitCanceled)
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
		err := analyze(ctx, kfsCore)
		if err == nil {
			setTaskStatus(StatusFinished)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			setTaskStatus(StatusCanceled)
			return
		}
		setTaskStatus(StatusError)
	}()
}

func analyze(ctx context.Context, kfsCore *core.KFS) error {
	setTaskStatus(StatusRunningCollect)
	hashList, err := kfsCore.Db.ListExpectFileType(ctx)
	if err != nil {
		return err
	}
	setTaskTotal(len(hashList))
	for _, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		ft, err := server.AnalyzeFileType(kfsCore, hash)
		if err != nil {
			addTaskError(err)
			continue
		}
		err = server.InsertExif(context.TODO(), kfsCore, hash, ft)
		if err != nil {
			addTaskError(err)
			continue
		}
		err = server.InsertFileType(context.TODO(), kfsCore, hash, ft)
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