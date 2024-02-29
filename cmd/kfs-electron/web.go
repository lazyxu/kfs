package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lazyxu/kfs/cmd/kfs-electron/task/local_file_filter"
	"github.com/lazyxu/kfs/rpc/client/local_file"
	"net"
	"net/http"
	"strconv"
)

func webServer(lis net.Listener) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	// Routes
	e.POST("/api/v1/startDriverLocalFile", apiStarDriverLocalFile)
	e.POST("/api/v1/startDriverLocalFileFilter", apiStarDriverLocalFileFilter)
	e.POST("/api/v1/startAllLocalFileSync", startAllLocalFileSync)
	e.GET("/api/v1/event/driverLocalFile/:driverId", local_file.ApiEvent)
	e.GET("/api/v1/event/driverLocalFileFilter/:driverId", local_file_filter.ApiEvent)

	// Start server
	e.Listener = lis
	e.Logger.Fatal(e.Start(""))
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func ok(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{Data: data})
}

// Handler

func apiStarDriverLocalFile(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	serverAddr := c.QueryParam("serverAddr")
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	srcPath := c.QueryParam("srcPath")
	ignores := c.QueryParam("ignores")
	encoder := c.QueryParam("encoder")
	ctx := c.Request().Context()
	d, err := local_file.GetOrLoadDriver(driverId)
	if err != nil {
		return err
	}
	d.StartOrStop(ctx, start, serverAddr, srcPath, ignores, encoder)
	return c.String(http.StatusOK, "")
}

func apiStarDriverLocalFileFilter(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	srcPath := c.QueryParam("srcPath")
	ignores := c.QueryParam("ignores")
	ctx := c.Request().Context()
	d, err := local_file_filter.GetOrLoadDriver(driverId)
	if err != nil {
		return err
	}
	d.StartOrStop(ctx, start, srcPath, ignores)
	return c.String(http.StatusOK, "")
}
