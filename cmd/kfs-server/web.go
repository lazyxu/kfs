package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
)

func webServer(webPortString string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Routes
	e.GET("/api/v1/branches", apiBranches)
	e.GET("/api/v1/open", apiOpen)

	println("KFS web server listening at:", webPortString)
	// Start server
	e.Logger.Fatal(e.Start(":" + webPortString))
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
func apiBranches(c echo.Context) error {
	branches, err := kfsCore.BranchList(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, branches)
}

func apiOpen(c echo.Context) error {
	branchName := c.QueryParam("branchName")
	filePath := c.QueryParam("filePath")
	maxContentSizeStr := c.QueryParam("maxContentSize")
	maxContentSize, err := strconv.ParseInt(maxContentSizeStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "maxContentSize should be a number")
	}
	dirItemOpened, err := kfsCore.Open2(c.Request().Context(), branchName, filePath, maxContentSize)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, dirItemOpened)
}
