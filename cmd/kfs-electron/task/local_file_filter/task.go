package local_file_filter

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
	Clients: make(map[echo.Context]core.EventClient[TaskInfo]),
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
	StatusWaitCanceled = 5
	StatusRunning      = 6
)

type TaskInfoClearable struct {
	Size      uint64 `json:"size"`
	FileCount uint64 `json:"fileCount"`
	DirCount  uint64 `json:"dirCount"`

	TotalSize      uint64 `json:"totalSize"`
	TotalFileCount uint64 `json:"totalFileCount"`
	TotalDirCount  uint64 `json:"totalDirCount"`
}

type TaskInfo struct {
	TaskInfoClearable

	cancel       context.CancelFunc
	Status       int   `json:"status"`
	LastDoneTime int64 `json:"lastDoneTime"`

	Cost int64 `json:"cost"`

	ErrMsg        string   `json:"errMsg"`
	Warnings      []string `json:"warnings"`
	startTime     time.Time
	CurFile       string   `json:"curFile"`
	CurSize       uint64   `json:"curSize"`
	CurDirItemCnt uint64   `json:"curDirItemCnt"`
	CurDir        string   `json:"curDir"`
	Ignores       []string `json:"ignores"`
}

func (d *DriverLocalFile) setTaskStatus(status int) {
	d.mutex.Lock()
	d.setTaskStatusWithLock(status)
	d.mutex.Unlock()
}

func (d *DriverLocalFile) setTaskError(errMsg string) {
	d.mutex.Lock()
	d.taskInfo.Status = StatusError
	d.taskInfo.cancel = nil
	d.taskInfo.LastDoneTime = time.Now().UnixNano()
	d.taskInfo.ErrMsg = errMsg
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) addTaskWarning(errMsg string) {
	d.mutex.Lock()
	d.taskInfo.Warnings = append(d.taskInfo.Warnings, errMsg)
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) addTaskIgnores(ignore string) {
	d.mutex.Lock()
	d.taskInfo.Ignores = append(d.taskInfo.Ignores, ignore)
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) setTaskStatusWithLock(status int) {
	d.taskInfo.Status = status
	if status == StatusRunning {
		d.taskInfo.Size = 0
		d.taskInfo.FileCount = 0
		d.taskInfo.DirCount = 0
		d.taskInfo.TotalSize = 0
		d.taskInfo.TotalFileCount = 0
		d.taskInfo.TotalDirCount = 0
		d.taskInfo.ErrMsg = ""
		d.taskInfo.CurFile = ""
		d.taskInfo.CurSize = 0
		d.taskInfo.CurDir = ""
		d.taskInfo.CurDirItemCnt = 0
		d.taskInfo.Warnings = make([]string, 0)
		d.taskInfo.Ignores = make([]string, 0)
		d.taskInfo.startTime = time.Now()
	} else if status == StatusFinished || status == StatusCanceled || status == StatusError {
		d.taskInfo.cancel = nil
		d.taskInfo.LastDoneTime = time.Now().UnixNano()
		d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	}
	s.SendAll()
}

func (d *DriverLocalFile) setTaskFile(path string, info os.FileInfo) {
	d.mutex.Lock()
	d.taskInfo.CurFile = path
	d.taskInfo.CurSize = uint64(info.Size())
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) setTaskDir(path string, n uint64) {
	d.mutex.Lock()
	d.taskInfo.CurDir = path
	d.taskInfo.CurDirItemCnt = n
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) addTaskTotal(info os.FileInfo) {
	d.mutex.Lock()
	if info.IsDir() {
		d.taskInfo.TotalDirCount++
	} else {
		d.taskInfo.TotalFileCount++
		d.taskInfo.TotalSize += uint64(info.Size())
	}
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

func (d *DriverLocalFile) addTaskCnt(info os.FileInfo) {
	d.mutex.Lock()
	d.taskInfo.CurFile = ""
	d.taskInfo.CurSize = 0
	d.taskInfo.CurDir = ""
	d.taskInfo.CurDirItemCnt = 0
	if info.IsDir() {
		d.taskInfo.DirCount++
	} else {
		d.taskInfo.FileCount++
		d.taskInfo.Size += uint64(info.Size())
	}
	d.taskInfo.Cost = time.Now().Sub(d.taskInfo.startTime).Milliseconds()
	s.SendAll()
	d.mutex.Unlock()
}

var drivers sync.Map

func GetOrLoadDriver(driverId uint64) (*DriverLocalFile, error) {
	d, _ := drivers.LoadOrStore(driverId, &DriverLocalFile{
		driverId: driverId,
		taskInfo: TaskInfo{
			cancel:   nil,
			Status:   StatusIdle,
			Warnings: make([]string, 0),
			Ignores:  make([]string, 0),
		},
		mutex: &sync.Mutex{},
	})
	return d.(*DriverLocalFile), nil
}

func (d *DriverLocalFile) StartOrStop(ctx context.Context, start bool, srcPath, ignores string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !start {
		d.setTaskStatusWithLock(StatusWaitCanceled)
		d.taskInfo.cancel()
		return
	}
	if d.taskInfo.Status == StatusRunning || d.taskInfo.Status == StatusWaitCanceled {
		return
	}
	d.setTaskStatusWithLock(StatusRunning)
	ctx, cancel := context.WithCancel(context.TODO())
	d.taskInfo.cancel = cancel
	go func() {
		err := d.DoFilter(ctx, srcPath, ignores)
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
