package main

import (
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/baidu_photo"
	"net/http"
	"strconv"
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

func apiUpdateDriverSync(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
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
	sStr := c.QueryParam("s")
	s, err := strconv.ParseInt(sStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "s should be a number")
	}
	err = kfsCore.Db.UpdateDriverSync(c.Request().Context(), id, sync, h, m, s)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiDeleteDriver(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	err = kfsCore.Db.DeleteDriver(c.Request().Context(), driverId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}
