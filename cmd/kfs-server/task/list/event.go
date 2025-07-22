package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
	"github.com/lazyxu/kfs/rpc/server"
)

type Response struct {
	Files  []*sqlc.ListLogicFilesByParentIdRow `json:"files,omitempty"`
	ErrMsg string                              `json:"errMsg,omitempty"`
	N      int                                 `json:"n,omitempty"`
	Dir    *sqlc.LogicFile                     `json:"dir,omitempty"`
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

	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	var userID int64 = 1
	dir, files, err := kfsCore.Q.ListDir(c.Request().Context(), userID, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}

	n := len(files)
	send(c, Response{Dir: dir, N: n})

	tick := time.NewTicker(time.Second)
	curFiles := []*sqlc.ListLogicFilesByParentIdRow{}
	packetSize := 1
	for _, file := range files {
		if !os.FileMode(file.Mode).IsDir() {
			err = server.AnalyzeIfNoFileType(c.Request().Context(), kfsCore, file.FileID.Int64, file.Size, file.Md5, file.Sha256)
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
			curFiles = []*sqlc.ListLogicFilesByParentIdRow{}
		default:
			if len(curFiles) > packetSize {
				send(c, Response{Files: curFiles, N: n})
				curFiles = []*sqlc.ListLogicFilesByParentIdRow{}
				packetSize = packetSize * 2
			}
		}
	}
	if len(curFiles) != 0 {
		send(c, Response{Files: curFiles, N: n})
	}
	return nil
}
