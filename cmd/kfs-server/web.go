package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func webServer(webPortString string) {
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
	e.GET("/api/v1/branches", apiBranches)
	e.GET("/api/v1/list", apiList)
	e.GET("/api/v1/openFile", apiOpenFile)
	e.GET("/api/v1/downloadFile", apiDownloadFile)

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

func apiList(c echo.Context) error {
	branchName := c.QueryParam("branchName")
	filePath := c.QueryParam("filePath")
	dirItems, err := kfsCore.List(c.Request().Context(), branchName, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, dirItems)
}

func apiOpenFile(c echo.Context) error {
	branchName := c.QueryParam("branchName")
	filePath := c.QueryParam("filePath")
	maxContentSizeStr := c.QueryParam("maxContentSize")
	maxContentSize, err := strconv.ParseInt(maxContentSizeStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "maxContentSize should be a number")
	}
	rc, tooLarge, err := kfsCore.OpenFile(c.Request().Context(), branchName, filePath, maxContentSize)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	if tooLarge {
		c.Response().Header().Set("Kfs-Too-Large", "true")
		return c.String(http.StatusOK, "")
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}

func apiDownloadFile(c echo.Context) error {
	branchName := c.QueryParam("branchName")
	filePath := c.QueryParam("filePath")
	rc, _, err := kfsCore.OpenFile(c.Request().Context(), branchName, filePath, -1)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}
