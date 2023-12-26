package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	jpegimage "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"image"
	"strings"
	"time"
)

func GetJpegExifData(kfsCore *core.KFS, hash string) (e dao.Exif, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	parse, err := jpegimage.NewJpegMediaParser().Parse(rc, int(rc.Size()))
	if err != nil {
		return
	}
	ifd, i, err := parse.Exif()
	if err != nil {
		return
	}
	thumbnail, err := ifd.Thumbnail()
	if err != nil && !errors.Is(err, exif.ErrNoThumbnail) {
		return
	}
	fmt.Printf("thumbnail: %d\n", len(thumbnail))
	ets, _, err := exif.GetFlatExifData(i, nil)
	if err != nil {
		return
	}
	for _, et := range ets {
		fmt.Printf("exif: %s %v\n", et.TagName, et.Value)
	}
	return
}

func getImageHeightWidth(kfsCore *core.KFS, hash string) (hw dao.HeightWidth, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	conf, _, err := image.DecodeConfig(rc)
	if err != nil {
		return
	}
	return dao.HeightWidth{
		Width:  uint64(conf.Width),
		Height: uint64(conf.Height),
	}, nil
}

const defaultOffset = "+08:00"

func insertImageTime(ctx context.Context, kfsCore *core.KFS, hash string, e dao.Exif) error {
	var t int64
	if e.DateTime != "" {
		offset := e.OffsetTime
		if offset == "" {
			offset = defaultOffset
		}
		tt, err := time.Parse("2006:01:02 15:04:05 -07:00", e.DateTime+" "+offset)
		if err != nil {
			return err
		}
		t = tt.UnixNano()
	} else if e.DateTimeOriginal != "" {
		offset := e.OffsetTime
		if offset == "" {
			offset = defaultOffset
		}
		tt, err := time.Parse("2006:01:02 15:04:05 -07:00", e.DateTimeOriginal+" "+offset)
		if err != nil {
			return err
		}
		t = tt.UnixNano()
	} else if e.DateTimeDigitized != "" {
		offset := e.OffsetTime
		if offset == "" {
			offset = defaultOffset
		}
		tt, err := time.Parse("2006:01:02 15:04:05 -07:00", e.DateTimeDigitized+" "+offset)
		if err != nil {
			return err
		}
		t = tt.UnixNano()
	}
	if t == 0 {
		t = kfsCore.Db.GetEarliestCrated(ctx, hash)
	}
	_, err := kfsCore.Db.InsertDCIMMetadataTime(ctx, hash, t)
	if err != nil {
		return err
	}
	return nil
}

func insertVideoTime(ctx context.Context, kfsCore *core.KFS, hash string, m dao.VideoMetadata) error {
	t := m.Created
	if t == 0 {
		t = m.Modified
	}
	if t == 0 {
		t = kfsCore.Db.GetEarliestCrated(ctx, hash)
	}
	_, err := kfsCore.Db.InsertDCIMMetadataTime(ctx, hash, t)
	if err != nil {
		return err
	}
	return nil
}

func InsertExif(ctx context.Context, kfsCore *core.KFS, hash string, fileType dao.FileType) error {
	if fileType.Type == "image" {
		//if fileType.SubType == matchers.TypeJpeg.MIME.Subtype {
		//	GetJpegExifData(kfsCore, hash)
		//}
		hw, err := getImageHeightWidth(kfsCore, hash)
		if err != nil {
			return err
		}
		_, err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
		// TODO: what if exist
		if err != nil {
			return err
		}
		e, err := GetExifData(kfsCore, hash)
		if err != nil {
			_, err = kfsCore.Db.InsertNullExif(ctx, hash)
			// TODO: what if exist
			if err != nil {
				return err
			}
			return nil
		}
		err = insertImageTime(ctx, kfsCore, hash, e)
		// TODO: what if exist
		if err != nil {
			return err
		}
		_, err = kfsCore.Db.InsertExif(ctx, hash, e)
		// TODO: what if exist
		if err != nil {
			return err
		}
		return nil
	} else if fileType.Type == "video" {
		m, hw, err := GetVideoMetadata(kfsCore, hash)
		if err != nil {
			_, err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
			// TODO: what if exist
			if err != nil {
				return err
			}
			_, err = kfsCore.Db.InsertNullVideoMetadata(ctx, hash)
			// TODO: what if exist
			if err != nil {
				return err
			}
			return nil
		}
		_, err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
		// TODO: what if exist
		if err != nil {
			return err
		}
		err = insertVideoTime(ctx, kfsCore, hash, m)
		// TODO: what if exist
		if err != nil {
			return err
		}
		_, err = kfsCore.Db.InsertVideoMetadata(ctx, hash, m)
		// TODO: what if exist
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func GetVideoMetadata(kfsCore *core.KFS, hash string) (m dao.VideoMetadata, hw dao.HeightWidth, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	var fileInfo VideoFile
	err = fileInfo.Open(rc)
	if err != nil {
		return
	}
	err = fileInfo.Parse()
	if err != nil {
		return
	}
	m = dao.VideoMetadata{
		Codec:    fileInfo.Codec,
		Created:  fileInfo.Movie.Created.UnixNano(),
		Modified: fileInfo.Movie.Modified.UnixNano(),
		Duration: fileInfo.Movie.Duration,
	}
	for _, track := range fileInfo.Movie.Tracks {
		if track.Height != 0 && track.Width != 0 {
			hw = dao.HeightWidth{
				Width:  uint64(track.Width),
				Height: uint64(track.Height),
			}
		}
	}
	return
}

func GetExifData(kfsCore *core.KFS, hash string) (e dao.Exif, err error) {
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
