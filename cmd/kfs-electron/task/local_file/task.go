package local_file

import (
	"context"
	"errors"
	"os"
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
		//startStr := c.QueryParam("start")
		//serverAddr := c.QueryParam("serverAddr")
		d, err := GetOrLoadDriver(driverId)
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
	Status       int   `json:"status"`
	LastDoneTime int64 `json:"lastDoneTime"`

	Size      uint64 `json:"size"`
	FileCount uint64 `json:"fileCount"`
	DirCount  uint64 `json:"dirCount"`

	TotalSize      uint64 `json:"totalSize"`
	TotalFileCount uint64 `json:"totalFileCount"`
	TotalDirCount  uint64 `json:"totalDirCount"`

	Cost int64 `json:"cost"`

	ErrMsg   string   `json:"errMsg"`
	Warnings []string `json:"warnings"`
}

func (d *DriverLocalFile) setTaskStatus(status int) {
	d.mutex.Lock()
	d.taskInfo.Status = status
	if status == StatusFinished || status == StatusCanceled || status == StatusError {
		d.taskInfo.cancel = nil
		d.taskInfo.LastDoneTime = time.Now().UnixNano()
	}
	if status == StatusWaitRunning || status == StatusRunning {
		d.taskInfo.Size = 0
		d.taskInfo.FileCount = 0
		d.taskInfo.DirCount = 0
		d.taskInfo.TotalSize = 0
		d.taskInfo.TotalFileCount = 0
		d.taskInfo.TotalDirCount = 0
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

func (d *DriverLocalFile) addTaskWarning(errMsg string) {
	d.mutex.Lock()
	d.taskInfo.Warnings = append(d.taskInfo.Warnings, errMsg)
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverLocalFile) setTaskCost(cost int64) {
	d.mutex.Lock()
	d.taskInfo.Cost = cost
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
		d.taskInfo.Size = 0
		d.taskInfo.FileCount = 0
		d.taskInfo.DirCount = 0
		d.taskInfo.TotalSize = 0
		d.taskInfo.TotalFileCount = 0
		d.taskInfo.TotalDirCount = 0
		d.taskInfo.Warnings = make([]string, 0)
	}
	s.SendAll()
}

func (d *DriverLocalFile) addTaskTotal(info os.FileInfo) {
	d.mutex.Lock()
	if info.IsDir() {
		d.taskInfo.TotalDirCount++
	} else {
		d.taskInfo.TotalFileCount++
		d.taskInfo.TotalSize += uint64(info.Size())
	}
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverLocalFile) addTaskCnt(info os.FileInfo) {
	d.mutex.Lock()
	if info.IsDir() {
		d.taskInfo.DirCount++
	} else {
		d.taskInfo.FileCount++
		d.taskInfo.Size += uint64(info.Size())
	}
	d.mutex.Unlock()
	s.SendAll()
}

var drivers sync.Map

func GetOrLoadDriver(driverId uint64) (*DriverLocalFile, error) {
	d, _ := drivers.LoadOrStore(driverId, &DriverLocalFile{
		driverId: driverId,
		taskInfo: TaskInfo{
			cancel:   nil,
			Status:   StatusIdle,
			Warnings: make([]string, 0),
		},
		mutex: &sync.Mutex{},
	})
	return d.(*DriverLocalFile), nil
}

func (d *DriverLocalFile) StartOrStop(ctx context.Context, start bool, serverAddr string, srcPath, encoder string) {
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
		err := d.Analyze(ctx, serverAddr, srcPath, encoder)
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

func (d *DriverLocalFile) Analyze(ctx context.Context, serverAddr string, srcPath, encoder string) error {
	d.setTaskStatus(StatusRunning)
	err := d.eventSourceBackup(ctx, d.driverId, srcPath, serverAddr, encoder)
	if err != nil {
		return err
	}
	return nil
}
