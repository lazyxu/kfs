package main

import (
	"context"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	"github.com/labstack/echo/v4"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var exifChan = make(chan bool)

var remain atomic.Uint64

func apiAnalysisExif(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	exifChan <- start
	return c.String(http.StatusOK, "")
}

func AnalysisExifProcess() {
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case start := <-exifChan:
				if start {
					if remain.Load() == 0 {
						go AnalysisExif(ctx, kfsCore)
					}
				} else {
					cancel()
				}
			}
		}
	}()
}

func AnalysisExif(ctx context.Context, fs *core.KFS) error {
	hashList, err := fs.Db.ListExpectExif(ctx)
	if err != nil {
		return err
	}
	remain.Store(uint64(len(hashList)))
	defer func() {
		remain.Store(0)
	}()
	for i, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		rc, err := fs.S.ReadWithSize(hash)
		if err != nil {
			return err
		}
		defer rc.Close()
		d, err := GetExifData(rc)
		if err != nil {
			_, err = fs.Db.InsertNullExif(ctx, hash)
			// TODO: what if exist
			if err != nil {
				println("InsertNullExif", err)
				return err
			}
			continue
		}
		fmt.Printf("%+v\n", d)
		_, err = fs.Db.InsertExif(ctx, hash, d)
		// TODO: what if exist
		if err != nil {
			println("InsertExif", err)
			return err
		}
		remain.Store(uint64(len(hashList) - i - 1))
	}
	if err != nil {
		return err
	}
	return nil
}

func GetExifData(r io.Reader) (d dao.ExifData, err error) {
	dt, err := exif.SearchAndExtractExifWithReader(r)
	if err != nil {
		return
	}
	ets, _, err := exif.GetFlatExifData(dt, nil)
	if err != nil {
		return
	}
	for _, et := range ets {
		if et.TagName == "ExifVersion" {
			d.Version = et.Value.(exifundefined.Tag9000ExifVersion).ExifVersion
		} else if et.TagName == "DateTime" {
			t, err := time.Parse(time.DateTime, et.Value.(string))
			if err != nil {
				continue
			}
			d.DateTime = uint64(t.UnixNano())
		} else if et.TagName == "HostComputer" {
			d.HostComputer = et.Value.(string)
		} else if et.TagName == "OffsetTime" {
			d.OffsetTime = et.Value.(string)
		} else if et.TagName == "HostComputer" {
			d.HostComputer = et.Value.(string)
		} else if et.TagName == "GPSLatitudeRef" {
			d.GPSLatitudeRef = et.Value.(string)
		} else if et.TagName == "GPSLatitude" {
			d.GPSLatitude = GPS2Float(et.Value.([]exifcommon.Rational))
		} else if et.TagName == "GPSLongitudeRef" {
			d.GPSLongitudeRef = et.Value.(string)
		} else if et.TagName == "GPSLongitude" {
			d.GPSLongitude = GPS2Float(et.Value.([]exifcommon.Rational))
		}
	}
	return d, nil
}

func GPS2Float(rational []exifcommon.Rational) float64 {
	if len(rational) == 3 {
		return float64(rational[0].Numerator)/float64(rational[0].Denominator) +
			float64(rational[1].Numerator)/float64(rational[1].Denominator)/60.0 +
			float64(rational[2].Numerator)/float64(rational[2].Denominator)/3600.0
	}
	return 0
}
