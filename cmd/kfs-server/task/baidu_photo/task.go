package baidu_photo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

type Client struct {
	ch chan TaskInfo
	d  *DriverBaiduPhoto
}

func (c *Client) Chan() chan TaskInfo {
	return c.ch
}

func (c *Client) Message() TaskInfo {
	fmt.Printf("Message ptr:%p id:%d, %+v\n", c.d, c.d.id, c.d.taskInfo)
	return c.d.taskInfo
}

var s = &core.EventServer[TaskInfo]{
	Clients: make(map[echo.Context]core.EventClient[TaskInfo]),
	NewClient: func(c echo.Context, kfsCore *core.KFS) (core.EventClient[TaskInfo], error) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 0)
		if err != nil {
			return nil, err
		}
		d, err := GetOrLoad(c.Request().Context(), kfsCore, id)
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
	fmt.Printf("setTaskStatus ptr:%p id:%d, %d\n", d, d.id, status)
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
	fmt.Printf("setTaskError ptr:%p id:%d, %s\n", d, d.id, errMsg)
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
	fmt.Printf("addTaskTotal ptr:%p id:%d, %d\n", d, d.id, total)
	d.mutex.Lock()
	d.taskInfo.Total += total
	d.mutex.Unlock()
	s.SendAll()
}

func (d *DriverBaiduPhoto) addTaskCnt() {
	fmt.Printf("addTaskCnt ptr:%p id:%d, %d\n", d, d.id, d.taskInfo.Cnt)
	d.mutex.Lock()
	d.taskInfo.Cnt++
	d.mutex.Unlock()
	s.SendAll()
}

var tasks sync.Map

func GetOrLoad(ctx context.Context, kfsCore *core.KFS, id int64) (*DriverBaiduPhoto, error) {
	d, ok := tasks.Load(id)
	if ok {
		return d.(*DriverBaiduPhoto), nil
	}
	t, err := LoadTaskFromDb(ctx, kfsCore, id)
	if err != nil {
		return nil, err
	}
	tasks.Store(id, t)
	return t, nil
}

func LoadTaskFromDb(ctx context.Context, kfsCore *core.KFS, id int64) (*DriverBaiduPhoto, error) {
	s, err := kfsCore.Q.GetSyncFromById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &DriverBaiduPhoto{
		kfsCore:    kfsCore,
		id:         id,
		BdussBfess: s.BdussBfess,
		Stoken:     s.Stoken,
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

func (d *DriverBaiduPhoto) StartOrStop2(ctx context.Context, start bool) {
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
		err := d.Analyze2(ctx)
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

func toStringSlice(s []BaiduFile) []string {
	c := make([]string, len(s))
	for i, v := range s {
		c[i] = v.Md5
	}
	return c
}

func (d *DriverBaiduPhoto) Analyze2(ctx context.Context) (err error) {
	d.setTaskStatus(StatusRunning)
	defer func() {
		if err != nil {
			return
		}
		_, err = d.kfsCore.Q.SetLivpForMovAndHeicOrJpgInDriver(ctx)
	}()
	var err1 error
	err = d.GetAllFile(ctx, func(list []BaiduFile) bool {
		d.addTaskTotal(len(list))
		var filterdList map[string]*sqlc.File
		filterdList, err1 = d.kfsCore.Q.ListFileMd5(ctx, toStringSlice(list))
		if err1 != nil {
			return false
		}
		for i, f := range list {
			fmt.Printf("[%d/%d] handle %s\n", i+1, len(list), f.Path)
			select {
			case <-ctx.Done():
				err1 = context.Canceled
				return false
			default:
			}
			err1 = d.Download(ctx, f, filterdList[f.Md5])
			if err1 != nil {
				return false
			}
			d.addTaskCnt()
		}
		return true
	})
	if err != nil {
		return
	}
	if err1 != nil {
		err = err1
		return
	}
	return
}
