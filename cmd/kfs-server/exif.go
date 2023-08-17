package main

import (
	"context"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/dao"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
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
	exifMap, err := kfsCore.Db.ListExif(c.Request().Context())
	if err != nil {
		return err
	}
	return ok(c, exifMap)
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
	for i, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		var e dao.Exif
		e, err = GetExifData(hash)
		if err != nil {
			fmt.Printf("%d %s NullExif\n", len(hashList)-i, hash)
			_, err = kfsCore.Db.InsertNullExif(ctx, hash)
			// TODO: what if exist
			if err != nil {
				println("InsertNullExif", err.Error())
				return err
			}
			exifCnt.Add(1)
			continue
		}
		fmt.Printf("%d %s %+v\n", len(hashList)-i, hash, e)
		if e.DateTime == 0 {
			println("d.DateTime == 0")
		}
		_, err = kfsCore.Db.InsertExif(ctx, hash, e)
		// TODO: what if exist
		if err != nil {
			println("InsertExif", err.Error())
			return err
		}
		exifCnt.Add(1)
	}
	if err != nil {
		return err
	}
	return nil
}

func GetExifData(hash string) (e dao.Exif, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	dt, err := exif.SearchAndExtractExifWithReader(rc)
	if err != nil {
		return
	}
	ets, _, err := exif.GetFlatExifData(dt, nil)
	if err != nil {
		return
	}
	for _, et := range ets {
		//fmt.Printf("%s %v\n", et.TagName, et.Value)
		if et.TagName == "ExifVersion" {
			e.Version = et.Value.(exifundefined.Tag9000ExifVersion).ExifVersion
		} else if et.TagName == "DateTime" || et.TagName == "DateTimeOriginal" {
			// TODO: time zone
			t, err := time.Parse("2006:01:02 15:04:05", et.Value.(string))
			if err != nil {
				println("time.Parse", et.Value.(string), err.Error())
				return e, err
			}
			e.DateTime = uint64(t.UnixNano())
		} else if et.TagName == "HostComputer" {
			e.HostComputer = et.Value.(string)
		} else if et.TagName == "OffsetTime" {
			e.OffsetTime = et.Value.(string)
		} else if et.TagName == "HostComputer" {
			e.HostComputer = et.Value.(string)
		} else if et.TagName == "GPSLatitudeRef" {
			e.GPSLatitudeRef = et.Value.(string)
		} else if et.TagName == "GPSLatitude" {
			e.GPSLatitude = GPS2Float(et.Value.([]exifcommon.Rational))
		} else if et.TagName == "GPSLongitudeRef" {
			e.GPSLongitudeRef = et.Value.(string)
		} else if et.TagName == "GPSLongitude" {
			e.GPSLongitude = GPS2Float(et.Value.([]exifcommon.Rational))
		}
	}
	return e, nil
}

func GPS2Float(rational []exifcommon.Rational) float64 {
	if len(rational) == 3 {
		return float64(rational[0].Numerator)/float64(rational[0].Denominator) +
			float64(rational[1].Numerator)/float64(rational[1].Denominator)/60.0 +
			float64(rational[2].Numerator)/float64(rational[2].Denominator)/3600.0
	}
	return 0
}
