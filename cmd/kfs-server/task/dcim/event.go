package dcim

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

type Response struct {
	DCIMList []*sqlc.ListDCIMRow `json:"dcimList,omitempty"`
	ErrMsg   string              `json:"errMsg,omitempty"`
	N        int                 `json:"n,omitempty"`
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

	fmt.Println("New DCIM event connection established")

	defer func() {
		fmt.Println("Closing DCIM event connection")
	}()

	dcimList, err := kfsCore.Q.ListDCIM(c.Request().Context())
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}

	n := len(dcimList)
	send(c, Response{N: n})

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	curDCIMList := []*sqlc.ListDCIMRow{}
	packetSize := 1

	for _, dcimItem := range dcimList {
		// 分析元数据（如果尚未分析）
		// if dcimItem.Type == "image" || dcimItem.Type == "video" {
		// 	err = server.AnalyzeMetadataNotAnalyzed(c.Request().Context(), kfsCore, sqlc.File{
		// 		ID:     dcimItem.ID,
		// 		Size:   dcimItem.Size,
		// 		Md5:    dcimItem.Md5,
		// 		Sha256: dcimItem.Sha256,
		// 	})
		// 	if errors.Is(err, context.Canceled) {
		// 		fmt.Println("DCIM connection canceled")
		// 		return nil
		// 	}
		// 	if err != nil {
		// 		send(c, Response{ErrMsg: err.Error()})
		// 	}
		// }

		curDCIMList = append(curDCIMList, &dcimItem)

		select {
		case <-tick.C:
			send(c, Response{DCIMList: curDCIMList, N: n})
			curDCIMList = []*sqlc.ListDCIMRow{}
		default:
			if len(curDCIMList) > packetSize {
				send(c, Response{DCIMList: curDCIMList, N: n})
				curDCIMList = []*sqlc.ListDCIMRow{}
				packetSize = packetSize * 2
			}
		}
	}

	if len(curDCIMList) != 0 {
		send(c, Response{DCIMList: curDCIMList, N: n})
	}

	return nil
}
