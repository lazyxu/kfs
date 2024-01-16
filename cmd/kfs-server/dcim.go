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

func apiListDCIMSearchType(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDCIMSearchType(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiListDCIMSearchSuffix(c echo.Context) error {
	drivers, err := kfsCore.Db.ListDCIMSearchSuffix(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiSearchDCIM(c echo.Context) error {
	typeList := c.QueryParams()["typeList[]"]
	if typeList == nil {
		typeList = []string{}
	}
	suffixList := c.QueryParams()["suffixList[]"]
	if suffixList == nil {
		suffixList = []string{}
	}
	drivers, err := kfsCore.Db.SearchDCIM(c.Request().Context(), typeList, suffixList)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}
