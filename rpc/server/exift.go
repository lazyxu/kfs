package server

import (
	"context"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"strings"
)

func InsertExif(ctx context.Context, kfsCore *core.KFS, hash string) (err error) {
	var e dao.Exif
	e, err = GetExifData(kfsCore, hash)
	if err != nil {
		_, err = kfsCore.Db.InsertNullExif(ctx, hash)
		// TODO: what if exist
		if err != nil {
			return err
		}
		return err
	}
	_, err = kfsCore.Db.InsertExif(ctx, hash, e)
	// TODO: what if exist
	if err != nil {
		return err
	}
	return nil
}

func GetExifData(kfsCore *core.KFS, hash string) (e dao.Exif, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	return GetExifDataWithReadAtSeeker(rc)
}

// func GetExifDataWithReadAtSeeker(rc io.Reader) (e dao.Exif, err error) {
func GetExifDataWithReadAtSeeker(rc exif.ReadAtSeeker) (e dao.Exif, err error) {
	//rs, err := exif.SearchAndExtractExifWithReader(rc)
	rs, err := exif.SearchAndExtractExifFromReadAtSeeker(rc)
	if err != nil {
		return
	}
	//ets, _, err := exif.GetFlatExifData(rs, nil)
	ets, _, err := exif.GetFlatExifDataUniversalSearchWithReadSeeker(rs, nil, false)
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
			e.Make = strings.TrimRight(et.Value.(string), "\x00")
		case "Model":
			e.Model = strings.TrimRight(et.Value.(string), "\x00")
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
