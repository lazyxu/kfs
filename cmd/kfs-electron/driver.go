package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/db/sqlc"
	"github.com/lazyxu/kfs/rpc/client/local_file"
)

type Param struct {
	ServerAddr string                          `json:"serverAddr"`
	List       []sqlc.LocalFileSyncToLogicFile `json:"drivers"`
}

func startAllLocalFileSync(c echo.Context) error {
	var p Param
	err := c.Bind(&p)
	if err != nil {
		return err
	}
	for _, d := range p.List {
		if d.Enabled && d.DeviceID == local_file.DeviceId {
			local_file.StartLocalFileSync(d.ID, p.ServerAddr, d.Cron, d.Srcpath, d.Ignores, d.Encoder)
		}
	}
	return c.String(http.StatusOK, "")
}
