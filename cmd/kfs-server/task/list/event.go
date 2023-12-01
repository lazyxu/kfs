package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
	"github.com/lazyxu/kfs/rpc/server"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	File   *dao.DriverFile `json:"file"`
	ErrMsg string          `json:"errMsg"`
	N      int             `json:"n"`
}

func send(c echo.Context, msg Response) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("invalid msg: %+v\n", msg)
	}
	fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
	c.Response().Flush()
}

func Handle(c echo.Context, kfsCore *core.KFS) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	fmt.Println("New connection established")

	defer func() {
		fmt.Println("Closing connection")
	}()

	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	files, err := kfsCore.ListDriverFile(c.Request().Context(), driverId, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}

	send(c, Response{N: len(files)})

	for _, file := range files {
		if !os.FileMode(file.Mode).IsDir() {
			err = TryAnalyze(c.Request().Context(), kfsCore, file.Hash)
			if errors.Is(err, context.Canceled) {
				fmt.Println("Connection closed")
				return nil
			}
			if err != nil {
				send(c, Response{File: &file, ErrMsg: err.Error()})
				continue
			}
		}
		send(c, Response{File: &file})
	}
	return nil
}

func TryAnalyze(ctx context.Context, kfsCore *core.KFS, hash string) error {
	ft, err := kfsCore.Db.GetFileType(ctx, hash)
	if errors.Is(err, dbBase.ErrNoRecords) {
		ft, err = server.AnalyzeFileType(kfsCore, hash)
		if err != nil {
			return err
		}
		err = server.InsertExif(ctx, kfsCore, hash, ft)
		if err != nil {
			return err
		}
		_, err = kfsCore.Db.InsertFileType(ctx, hash, ft)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}
