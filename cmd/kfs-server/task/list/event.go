package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/rpc/server"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
)

type Response struct {
	Files  []dao.DriverFile `json:"files,omitempty"`
	ErrMsg string           `json:"errMsg,omitempty"`
	N      int              `json:"n,omitempty"`
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
	files, err := kfsCore.Db.ListDriverFile(c.Request().Context(), driverId, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}

	n := len(files)
	send(c, Response{N: n})

	tick := time.NewTicker(time.Second)
	curFiles := []dao.DriverFile{}
	packetSize := 1
	for _, file := range files {
		if !os.FileMode(file.Mode).IsDir() {
			err = server.AnalyzeIfNoFileType(c.Request().Context(), kfsCore, file.Hash)
			if errors.Is(err, context.Canceled) {
				fmt.Println("Connection canceled")
				return nil
			}
			if err != nil {
				send(c, Response{ErrMsg: err.Error()})
			}
		}
		curFiles = append(curFiles, file)
		select {
		case <-tick.C:
			send(c, Response{Files: curFiles, N: n})
			curFiles = []dao.DriverFile{}
		default:
			if len(curFiles) > packetSize {
				send(c, Response{Files: curFiles, N: n})
				curFiles = []dao.DriverFile{}
				packetSize = packetSize * 2
			}
		}
	}
	if len(curFiles) != 0 {
		send(c, Response{Files: curFiles, N: n})
	}
	return nil
}
