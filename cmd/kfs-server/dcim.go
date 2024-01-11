package main

import "github.com/labstack/echo/v4"

func apiListDCIMDriver(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDCIMDriver(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiListDCIMMediaType(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDCIMMediaType(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiListDCIMLocation(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDCIMLocation(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}
