package main

import (
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/client/local_file"
	"net/http"
)

type Param struct {
	ServerAddr string       `json:"serverAddr"`
	Drivers    []dao.Driver `json:"drivers"`
}

func startAllLocalFileSync(c echo.Context) error {
	var p Param
	err := c.Bind(&p)
	if err != nil {
		return err
	}
	for _, d := range p.Drivers {
		local_file.StartLocalFileSync(d.Id, p.ServerAddr, d.H, d.M, d.SrcPath, d.Ignores, d.Encoder)
	}
	return c.String(http.StatusOK, "")
}
