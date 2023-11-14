package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/baidu_photo"
	"github.com/lazyxu/kfs/db/dbBase"
	"github.com/robfig/cron/v3"
	"net/http"
	"strconv"
	"sync"
)

func apiDrivers(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDriver(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiNewDriver(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	exist, err := kfsCore.Db.InsertDriver(c.Request().Context(), name, description)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, exist)
}

func apiNewDriverBaiduPhoto(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	code := c.QueryParam("code")
	exist, err := baidu_photo.InsertDriverBaiduPhoto(c.Request().Context(), kfsCore, name, description, code)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, exist)
}

func apiNewDriverLocalFile(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	deviceIdStr := c.QueryParam("deviceId")
	deviceId, err := strconv.ParseUint(deviceIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	srcPath := c.QueryParam("srcPath")
	encoder := c.QueryParam("encoder")
	concurrentStr := c.QueryParam("concurrent")
	concurrent, err := strconv.ParseUint(concurrentStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "concurrent should be a number")
	}
	exist, err := kfsCore.Db.InsertDriverLocalFile(c.Request().Context(), name, description, deviceId, srcPath, encoder, int(concurrent))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, exist)
}

func apiGetDriverSync(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	d, err := kfsCore.Db.GetDriverSync(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, d)
}

func apiGetDriverLocalFile(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	d, err := kfsCore.Db.GetDriverLocalFile(c.Request().Context(), driverId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, d)
}

func apiUpdateDriverSync(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	syncStr := c.QueryParam("sync")
	sync, err := strconv.ParseBool(syncStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "sync should be a boolean")
	}
	hStr := c.QueryParam("h")
	h, err := strconv.ParseInt(hStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "h should be a number")
	}
	mStr := c.QueryParam("m")
	m, err := strconv.ParseInt(mStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "m should be a number")
	}
	err = kfsCore.Db.UpdateDriverSync(c.Request().Context(), driverId, sync, h, m)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	d, err := kfsCore.Db.GetDriver(c.Request().Context(), driverId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	if d.Typ == dbBase.DRIVER_TYPE_BAIDU_PHOTO {
		startCloudSync(driverId, h, m)
	}
	return c.String(http.StatusOK, "")
}

var cronTasks sync.Map

type CronTask struct {
	c      *cron.Cron
	cancel context.CancelFunc
}

func startAllCloudSync() {
	drivers, err := kfsCore.Db.ListCloudDriverSync(context.TODO())
	if err != nil {
		panic(err)
	}
	for _, d := range drivers {
		startCloudSync(d.Id, d.H, d.M)
	}
}

func startCloudSync(driverId uint64, h int64, m int64) {
	actual, loaded := cronTasks.LoadOrStore(driverId, CronTask{
		c:      cron.New(),
		cancel: nil,
	})
	t := actual.(CronTask)
	if loaded {
		if t.cancel != nil {
			t.cancel()
			t.cancel = nil
		}
		t.c.Stop()
	}
	spec := fmt.Sprintf("%d %d * * ?", m, h)
	_, err := t.c.AddFunc(spec, func() {
		ctx, cancel := context.WithCancel(context.TODO())
		t.cancel = cancel
		d, err := baidu_photo.GetOrLoadDriver(ctx, kfsCore, driverId)
		if err != nil {
			cronTasks.LoadAndDelete(driverId)
			return
		}
		d.StartOrStop(ctx, true)
	})
	if err != nil {
		panic(err)
	}
	t.c.Start()
}

func apiDeleteDriver(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	err = kfsCore.Db.DeleteDriver(c.Request().Context(), driverId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}
