package main

import (
	"net/http"

	"github.com/lazyxu/kfs/db/sqlc"

	"github.com/labstack/echo/v4"
)

func apiNewDevice(c echo.Context) error {
	var p sqlc.InsertDeviceParams
	err := c.Bind(&p)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	err = kfsCore.Q.InsertDevice(c.Request().Context(), p)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiDeleteDevice(c echo.Context) error {
	deviceId := c.QueryParam("deviceId")
	err := kfsCore.Q.DeleteDevice(c.Request().Context(), deviceId)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiDevices(c echo.Context) error {
	devices, err := kfsCore.Q.ListDevice(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, devices)
}
