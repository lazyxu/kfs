package baidu_photo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/common"
	"github.com/lazyxu/kfs/core"
)

type Client struct {
	ch chan TaskInfo
	d  *DriverBaiduPhoto
}

func (c *Client) Chan() chan TaskInfo {
	return c.ch
}

func (c *Client) Message() TaskInfo {
	return c.d.taskInfo
}

var s = &common.EventServer[TaskInfo]{
	NewClient: func(c echo.Context, kfsCore *core.KFS) (common.Client[TaskInfo], error) {
		driverName := c.Param("name")
		d, err := GetOrLoadDriver(c.Request().Context(), kfsCore, driverName)
		if err != nil {
			return nil, err
		}
		return &Client{
			ch: make(chan TaskInfo),
			d:  d,
		}, nil
	},
}

func ApiEvent(c echo.Context, kfsCore *core.KFS) error {
	return s.Handle(c, kfsCore)
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

func (d *DriverBaiduPhoto) setTaskStatus(status int) {
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

func (d *DriverBaiduPhoto) setTaskError(errMsg string) {
	d.mutex.Lock()
	d.taskInfo.Status = StatusError
	d.taskInfo.cancel = nil
	d.taskInfo.LastDoneTime = time.Now().UnixNano()
	d.taskInfo.ErrMsg = errMsg
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverBaiduPhoto) setTaskStatusWithLock(status int) {
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

func (d *DriverBaiduPhoto) addTaskTotal(total int) {
	d.mutex.Lock()
	d.taskInfo.Total += total
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverBaiduPhoto) addTaskCnt() {
	d.mutex.Lock()
	d.taskInfo.Cnt++
	d.mutex.Unlock()
	s.SendAll()
}

var drivers sync.Map

func GetOrLoadDriver(ctx context.Context, kfsCore *core.KFS, driverName string) (*DriverBaiduPhoto, error) {
	d, ok := drivers.Load(driverName)
	if ok {
		return d.(*DriverBaiduPhoto), nil
	}
	driver, err := LoadDriverFromDb(ctx, kfsCore, driverName)
	if err != nil {
		return nil, err
	}
	drivers.Store(driverName, driver)
	return driver, nil
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

func (d *DriverBaiduPhoto) StartOrStop(ctx context.Context, start bool) {
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

func toStringSlice(s []File) []string {
	c := make([]string, len(s))
	for i, v := range s {
		c[i] = v.Md5
	}
	return c
}

func (d *DriverBaiduPhoto) Analyze(ctx context.Context) error {
	d.setTaskStatus(StatusRunning)
	var err1 error
	err := d.GetAllFile(ctx, func(list []File) bool {
		d.addTaskTotal(len(list))
		var m map[string]string
		m, err1 = d.kfsCore.Db.ListFileMd5(ctx, toStringSlice(list))
		if err1 != nil {
			return false
		}
		for i, f := range list {
			fmt.Printf("[%d/%d] handle %s\n", i, len(list), f.Path)
			select {
			case <-ctx.Done():
				err1 = context.Canceled
				return false
			default:
			}
			err1 = d.Download(ctx, f, m[f.Md5])
			if err1 != nil {
				return false
			}
			d.addTaskCnt()
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
