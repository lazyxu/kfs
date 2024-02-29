package main

import (
	"github.com/lazyxu/kfs/dao"
	"net/http"

	"github.com/labstack/echo/v4"
)

func apiNewDevice(c echo.Context) error {
	var d dao.Device
	err := c.Bind(&d)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	err = kfsCore.Db.InsertDevice(c.Request().Context(), d.Id, d.Name, d.OS, d.UserAgent, d.Hostname)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiDeleteDevice(c echo.Context) error {
	deviceId := c.QueryParam("deviceId")
	err := kfsCore.Db.DeleteDevice(c.Request().Context(), deviceId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiDevices(c echo.Context) error {
	devices, err := kfsCore.Db.ListDevice(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, devices)
}
