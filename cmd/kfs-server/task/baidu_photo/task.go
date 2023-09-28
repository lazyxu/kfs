package baidu_photo

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/common"
	"github.com/lazyxu/kfs/core"
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
	StatusIdle         = 0
	StatusFinished     = 1
	StatusCanceled     = 2
	StatusError        = 3
	StatusWaitRunning  = 4
	StatusWaitCanceled = 5
	StatusRunning      = 6
)

type TaskInfo struct {
	cancel       context.CancelFunc
	Status       int   `json:"status"`
	LastDoneTime int64 `json:"lastDoneTime"`
	Cnt          int   `json:"cnt"`
	Total        int   `json:"total"`
}

var taskInfo = TaskInfo{
	cancel: nil,
	Status: StatusIdle,
	Cnt:    0,
	Total:  0,
}

var mutex = &sync.RWMutex{}

func setTaskStatus(status int) {
	mutex.Lock()
	taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		taskInfo.cancel = nil
		taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunning {
		taskInfo.Cnt = 0
		taskInfo.Total = 0
	}
	mutex.Unlock()
	s.SendAll()
}

func addTaskTotal(total int) {
	mutex.Lock()
	taskInfo.Total += total
	mutex.Unlock()
	s.SendAll()
}

func addTaskCnt() {
	mutex.Lock()
	taskInfo.Cnt++
	mutex.Unlock()
	s.SendAll()
}

func LoadDriverFromDb(ctx context.Context, kfsCore *core.KFS, driverName string) (*DriverBaiduPhoto, error) {
	driver, err := kfsCore.Db.GetDriver(ctx, driverName)
	if err != nil {
		return nil, err
	}
	return &DriverBaiduPhoto{
		kfsCore:      kfsCore,
		driverName:   driverName,
		AccessToken:  driver.AccessToken,
		RefreshToken: driver.RefreshToken,
		AppKey:       AppKey,
		SecretKey:    SecretKey,
	}, nil
}

func StartOrStop(ctx context.Context, start bool, doTask func() error) {
	mutex.Lock()
	defer mutex.Unlock()
	if !start {
		setTaskStatus(StatusWaitCanceled)
		taskInfo.cancel()
		return
	}
	if taskInfo.Status == StatusWaitRunning ||
		taskInfo.Status == StatusRunning ||
		taskInfo.Status == StatusWaitCanceled {
		return
	}
	taskInfo.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	taskInfo.cancel = cancel
	go func() {
		err := doTask()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				setTaskStatus(StatusCanceled)
				return
			}
			setTaskStatus(StatusError)
			return
		}
		setTaskStatus(StatusFinished)
	}()
}

func (d *DriverBaiduPhoto) Analyze(ctx context.Context) error {
	setTaskStatus(StatusRunning)
	var err1 error
	err := d.GetAllFile(ctx, func(list []File) bool {
		addTaskTotal(len(list))
		for i, f := range list {
			fmt.Printf("[%d/%d] downloading %s\n", i, len(list), f.Path)
			select {
			case <-ctx.Done():
				err1 = context.DeadlineExceeded
				return false
			default:
			}
			err1 = d.Download(ctx, f)
			if err1 != nil {
				return false
			}
			addTaskCnt()
		}
		return true
	})
	if err != nil {
		return err
	}
	if err1 != nil {
		return err1
	}
	return nil
}
