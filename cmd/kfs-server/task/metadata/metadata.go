package metadata

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/common"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/server"
	"net/http"
	"sync"
	"time"
)

var s = common.EventServer[TaskInfo]{
	Message: func() TaskInfo {
		return taskInfo
	},
}

func ApiEventMetadataAnalysisTask(c echo.Context) error {
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

var metadataAnalysisTaskMutex = &sync.RWMutex{}

func setTaskStatus(status int) {
	metadataAnalysisTaskMutex.Lock()
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
	metadataAnalysisTaskMutex.Unlock()
	s.SendAll()
}

func setTaskTotal(total int) {
	metadataAnalysisTaskMutex.Lock()
	taskInfo.Status = StatusRunningAnalyze
	taskInfo.Errors = make([]string, 0)
	taskInfo.Cnt = 0
	taskInfo.Total = total
	metadataAnalysisTaskMutex.Unlock()
	s.SendAll()
}

func addTaskCnt() {
	metadataAnalysisTaskMutex.Lock()
	taskInfo.Cnt++
	metadataAnalysisTaskMutex.Unlock()
	s.SendAll()
}

func addTaskError(err error) {
	metadataAnalysisTaskMutex.Lock()
	taskInfo.Errors = append(taskInfo.Errors, err.Error())
	metadataAnalysisTaskMutex.Unlock()
	s.SendAll()
}

func StartMetadataAnalysisTask(c echo.Context, kfsCore *core.KFS, start bool) error {
	metadataAnalysisTaskMutex.Lock()
	defer metadataAnalysisTaskMutex.Unlock()
	if !start {
		setTaskStatus(StatusWaitCanceled)
		taskInfo.cancel()
		return c.String(http.StatusOK, "")
	}
	tryStartMetadataAnalysisTask(kfsCore)
	return c.String(http.StatusOK, "")
}

func tryStartMetadataAnalysisTask(kfsCore *core.KFS) {
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
		err := analyzeMetadata(ctx, kfsCore)
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

func analyzeMetadata(ctx context.Context, kfsCore *core.KFS) error {
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
