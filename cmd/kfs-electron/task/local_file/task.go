package local_file

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
)

type Client struct {
	ch chan TaskInfo
	d  *DriverLocalFile
}

func (c *Client) Chan() chan TaskInfo {
	return c.ch
}

func (c *Client) Message() TaskInfo {
	return c.d.taskInfo
}

var s = &core.EventServer[TaskInfo]{
	NewClient: func(c echo.Context, kfsCore *core.KFS) (core.EventClient[TaskInfo], error) {
		driverIdStr := c.Param("driverId")
		driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
		if err != nil {
			return nil, err
		}
		d, err := GetOrLoadDriver(c.Request().Context(), kfsCore, driverId)
		if err != nil {
			return nil, err
		}
		return &Client{
			ch: make(chan TaskInfo),
			d:  d,
		}, nil
	},
}

func ApiEvent(c echo.Context) error {
	return s.Handle(c, nil)
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
	Status       int    `json:"status"`
	LastDoneTime int64  `json:"lastDoneTime"`
	Cnt          int    `json:"cnt"`
	Total        int    `json:"total"`
	ErrMsg       string `json:"errMsg"`
}

func (d *DriverLocalFile) setTaskStatus(status int) {
	d.mutex.Lock()
	d.taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		d.taskInfo.cancel = nil
		d.taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunning {
		d.taskInfo.Cnt = 0
		d.taskInfo.Total = 0
		d.taskInfo.ErrMsg = ""
	}
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverLocalFile) setTaskError(errMsg string) {
	d.mutex.Lock()
	d.taskInfo.Status = StatusError
	d.taskInfo.cancel = nil
	d.taskInfo.LastDoneTime = time.Now().UnixNano()
	d.taskInfo.ErrMsg = errMsg
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverLocalFile) setTaskStatusWithLock(status int) {
	d.taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		d.taskInfo.cancel = nil
		d.taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunning {
		d.taskInfo.Cnt = 0
		d.taskInfo.Total = 0
	}
	s.SendAll()
}

func (d *DriverLocalFile) addTaskTotal(total int) {
	d.mutex.Lock()
	d.taskInfo.Total += total
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverLocalFile) addTaskCnt() {
	d.mutex.Lock()
	d.taskInfo.Cnt++
	d.mutex.Unlock()
	s.SendAll()
}

var drivers sync.Map

func GetOrLoadDriver(ctx context.Context, kfsCore *core.KFS, driverId uint64) (*DriverLocalFile, error) {
	d, ok := drivers.Load(driverId)
	if ok {
		return d.(*DriverLocalFile), nil
	}
	driver, err := LoadDriverFromDb(ctx, kfsCore, driverId)
	if err != nil {
		return nil, err
	}
	drivers.Store(driverId, driver)
	return driver, nil
}

func LoadDriverFromDb(ctx context.Context, kfsCore *core.KFS, driverId uint64) (*DriverLocalFile, error) {
	driver, err := kfsCore.Db.GetDriverToken(ctx, driverId)
	if err != nil {
		return nil, err
	}
	return &DriverLocalFile{
		kfsCore:      kfsCore,
		driverId:     driverId,
		AccessToken:  driver.AccessToken,
		RefreshToken: driver.RefreshToken,
		AppKey:       AppKey,
		SecretKey:    SecretKey,
		taskInfo: TaskInfo{
			cancel:       nil,
			Status:       StatusIdle,
			LastDoneTime: 0,
			Cnt:          0,
			Total:        0,
		},
		mutex: &sync.Mutex{},
	}, nil
}

func (d *DriverLocalFile) StartOrStop(ctx context.Context, start bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !start {
		d.setTaskStatusWithLock(StatusWaitCanceled)
		d.taskInfo.cancel()
		return
	}
	if d.taskInfo.Status == StatusWaitRunning ||
		d.taskInfo.Status == StatusRunning ||
		d.taskInfo.Status == StatusWaitCanceled {
		return
	}
	d.taskInfo.Status = StatusWaitRunning
	ctx, cancel := context.WithCancel(context.TODO())
	d.taskInfo.cancel = cancel
	go func() {
		err := d.Analyze(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				d.setTaskStatus(StatusCanceled)
				return
			}
			d.setTaskError(err.Error())
			return
		}
		d.setTaskStatus(StatusFinished)
	}()
}

func (d *DriverLocalFile) Analyze(ctx context.Context) error {
	d.setTaskStatus(StatusRunning)
	var err1 error
	err := d.eventSourceBackup(ctx, task.Name, task.Description, task.SrcPath, serverAddr, task.DriverId, task.DstPath, task.Encoder, task.Concurrent)
	if err != nil {
		return err
	}
	if err1 != nil {
		return err1
	}
	return nil
}
