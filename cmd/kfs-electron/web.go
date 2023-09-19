package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"net/http"
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
	e.GET("/api/v1/backupTask", apiBackupTask)
	e.GET("/api/v1/event/backupTask", apiEventBackupTask)
	e.POST("/api/v1/backupTask", apiNewBackupTask)
	e.DELETE("/api/v1/backupTask", apiDeleteBackupTask)
	e.POST("/api/v1/startBackupTask", apiStartBackupTask)
	e.GET("/api/v1/event/backupTask/:name", apiEventBackupTaskDetail)

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