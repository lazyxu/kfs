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
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	// Routes
	e.GET("/api/v1/backupTask", apiBackupTask)
	e.POST("/api/v1/backupTask", apiNewBackupTask)

	println("KFS electron web server listening at:", webPortString)
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

func apiBackupTask(c echo.Context) error {
	list, err := db.ListBackupTask(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, list)
}

func apiNewBackupTask(c echo.Context) error {
	name := c.QueryParam("name")
	description := c.QueryParam("description")
	srcPath := c.QueryParam("srcPath")
	driverName := c.QueryParam("driverName")
	dstPath := c.QueryParam("dstPath")
	encoder := c.QueryParam("encoder")
	concurrentStr := c.QueryParam("concurrent")
	concurrent, err := strconv.Atoi(concurrentStr)
	if err != nil {
		return err
	}
	err = upsertBackup(c.Request().Context(), db, name, description, srcPath, driverName, dstPath, encoder, concurrent)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "")
}
