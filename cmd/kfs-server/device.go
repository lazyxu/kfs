package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func apiNewDevice(c echo.Context) error {
	name := c.QueryParam("name")
	os := c.QueryParam("os")
	deviceId, err := kfsCore.Db.InsertDevice(c.Request().Context(), name, os)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, deviceId)
}

func apiDeleteDevice(c echo.Context) error {
	deviceIdStr := c.QueryParam("deviceId")
	deviceId, err := strconv.ParseUint(deviceIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	err = kfsCore.Db.DeleteDevice(c.Request().Context(), deviceId)
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
