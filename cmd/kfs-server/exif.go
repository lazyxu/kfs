package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/rpc/server"
	"net/http"
	"strconv"
	"sync/atomic"
)

var exifChan = make(chan bool)

var exifAnalyzing atomic.Bool
var exifFinished atomic.Bool
var exifCnt atomic.Uint64
var exifTotal atomic.Uint64

func apiAnalysisExif(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	exifChan <- start
	return c.String(http.StatusOK, "")
}

type ExifStatus struct {
	Analyzing bool   `json:"analyzing"`
	Finished  bool   `json:"finished"`
	Cnt       uint64 `json:"cnt"`
	Total     uint64 `json:"total"`
}

func apiExifStatus(c echo.Context) error {
	return ok(c, ExifStatus{
		Analyzing: exifAnalyzing.Load(),
		Finished:  exifFinished.Load(),
		Cnt:       exifCnt.Load(),
		Total:     exifTotal.Load(),
	})
}

func apiListExif(c echo.Context) error {
	data, err := kfsCore.Db.ListExifWithFileType(c.Request().Context())
	if err != nil {
		return err
	}
	return ok(c, data)
}

func apiGetMetadata(c echo.Context) error {
	hash := c.QueryParam("hash")
	data, err := kfsCore.Db.GetMetadata(c.Request().Context(), hash)
	if err != nil {
		return err
	}
	return ok(c, data)
}

func AnalysisExifProcess() {
	var ctx context.Context
	var cancel context.CancelFunc
	go func() {
		for {
			select {
			case start := <-exifChan:
				if start {
					if exifAnalyzing.CompareAndSwap(false, true) {
						ctx, cancel = context.WithCancel(context.TODO())
						go AnalysisExif(ctx)
					}
				} else {
					if cancel != nil {
						cancel()
					}
				}
			}
		}
	}()
}

func AnalysisExif(ctx context.Context) (err error) {
	println("AnalysisExif")
	exifFinished.Store(false)
	exifCnt.Store(0)
	exifTotal.Store(0)
	defer func() {
		exifAnalyzing.Store(false)
		exifFinished.Store(true)
	}()
	hashList, err := kfsCore.Db.ListExpectExif(ctx)
	if err != nil {
		return err
	}
	exifTotal.Store(uint64(len(hashList)))
	for _, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		ft, err := server.InsertFileType(ctx, kfsCore, hash)
		if err != nil {
			println("InsertFileType", err.Error())
			exifCnt.Add(1)
			continue
		}
		err = server.InsertExif(ctx, kfsCore, hash, ft)
		if err != nil {
			println("InsertExif", err.Error())
			exifCnt.Add(1)
			continue
		}
		exifCnt.Add(1)
	}
	if err != nil {
		return err
	}
	return nil
}
