package main

import (
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/labstack/echo/v4"
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

func apiListMetadata(c echo.Context) error {
	data, err := kfsCore.Db.ListMetadata(c.Request().Context())
	if err != nil {
		return err
	}
	return ok(c, data)
}

func apiListMetadataTime(c echo.Context) error {
	data, err := kfsCore.Db.ListMetadataTime(c.Request().Context())
	if err != nil {
		return err
	}
	return ok(c, data)
}

func apiGetMetadata(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	data, err := kfsCore.Q.GetMetadata(c.Request().Context(), id)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
	return ok(c, data)
}
