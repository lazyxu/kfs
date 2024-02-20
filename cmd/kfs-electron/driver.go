package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-electron/task/local_file"
	"github.com/lazyxu/kfs/dao"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
)

var cronTasks sync.Map

type CronTask struct {
	c      *cron.Cron
	id     cron.EntryID
	cancel context.CancelFunc
}

type Param struct {
	ServerAddr string       `json:"serverAddr"`
	Drivers    []dao.Driver `json:"drivers"`
}

func startAllLocalFileSync(c echo.Context) error {
	var p Param
	err := c.Bind(&p)
	if err != nil {
		return err
	}
	for _, d := range p.Drivers {
		startLocalFileSync(d.Id, p.ServerAddr, d.H, d.M, d.SrcPath, d.Ignores, d.Encoder)
	}
	return c.String(http.StatusOK, "")
}

func startLocalFileSync(driverId uint64, serverAddr string, h int64, m int64, srcPath, ignores, encoder string) {
	actual, loaded := cronTasks.LoadOrStore(driverId, &CronTask{
		c:      cron.New(),
		id:     -1,
		cancel: nil,
	})
	t := actual.(*CronTask)
	if loaded {
		if t.cancel != nil {
			t.cancel()
			t.cancel = nil
		}
		t.c.Remove(t.id)
	}
	spec := fmt.Sprintf("%d %d * * ?", m, h)
	var err error
	t.id, err = t.c.AddFunc(spec, func() {
		ctx, cancel := context.WithCancel(context.TODO())
		t.cancel = cancel
		d, err := local_file.GetOrLoadDriver(driverId)
		if err != nil {
			cronTasks.LoadAndDelete(driverId)
			return
		}
		d.StartOrStop(ctx, true, serverAddr, srcPath, ignores, encoder)
	})
	if err != nil {
		panic(err)
	}
	t.c.Start()
}
