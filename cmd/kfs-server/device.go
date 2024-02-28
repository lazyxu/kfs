package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func apiNewDevice(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")
	os := c.QueryParam("os")
	userAgent := c.QueryParam("userAgent")
	hostname := c.QueryParam("hostname")
	err := kfsCore.Db.InsertDevice(c.Request().Context(), id, name, os, userAgent, hostname)
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
