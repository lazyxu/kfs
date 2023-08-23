package main

import (
	"context"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	"github.com/h2non/filetype/types"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/dao"
	"net/http"
	"strconv"
	"strings"
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
		var fileType types.Type
		fileType, err = GetFileType(hash)
		if err != nil {
			println("GetFileType", err.Error())
			return err
		}
		_, err = kfsCore.Db.InsertFileType(ctx, hash, dao.FileType{
			Type:      fileType.MIME.Type,
			SubType:   fileType.MIME.Subtype,
			Extension: fileType.Extension,
		})
		if err != nil {
			println("InsertFileType", err.Error())
			return err
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
		switch et.TagName {
		case "ExifVersion":
			e.ExifVersion = et.Value.(exifundefined.Tag9000ExifVersion).ExifVersion
		case "ImageDescription":
			e.ImageDescription = et.Value.(string)
		case "Orientation":
			val := et.Value.([]uint16)
			if len(val) == 0 {
				e.Orientation = 0xFFFF
			} else if len(val) == 1 {
				e.Orientation = val[0]
			} else {
				panic(val)
			}
		case "DateTime":
			e.DateTime = et.Value.(string)
		case "DateTimeOriginal":
			e.DateTimeOriginal = et.Value.(string)
		case "DateTimeDigitized":
			e.DateTimeDigitized = et.Value.(string)
		case "OffsetTime":
			e.OffsetTime = et.Value.(string)
		case "OffsetTimeOriginal":
			e.OffsetTimeOriginal = et.Value.(string)
		case "OffsetTimeDigitized":
			e.OffsetTimeDigitized = et.Value.(string)
		case "SubsecTime":
			e.SubsecTime = et.Value.(string)
		case "SubsecTimeOriginal":
			e.SubsecTimeOriginal = et.Value.(string)
		case "SubsecTimeDigitized":
			e.SubsecTimeDigitized = et.Value.(string)
		case "HostComputer":
			e.HostComputer = et.Value.(string)
		case "Make":
			e.Make = strings.TrimSuffix(et.Value.(string), "\x00")
		case "Model":
			e.Model = et.Value.(string)
		case "ExifImageWidth":
			e.ExifImageWidth = et.Value.(uint64)
		case "ExifImageLength":
			e.ExifImageLength = et.Value.(uint64)
		case "GPSLatitudeRef":
			e.GPSLatitudeRef = et.Value.(string)
		case "GPSLatitude":
			e.GPSLatitude = GPS2Float(et.Value.([]exifcommon.Rational))
		case "GPSLongitudeRef":
			e.GPSLongitudeRef = et.Value.(string)
		case "GPSLongitude":
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
