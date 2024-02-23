package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/adrium/goheif"
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	exifundefined "github.com/dsoprea/go-exif/v3/undefined"
	jpegimage "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"image"
	"regexp"
	"strconv"
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
		var img image.Image
		img, err = goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
		if err != nil {
			return
		}
		return dao.HeightWidth{
			Width:  uint64(img.Bounds().Dx()),
			Height: uint64(img.Bounds().Dy()),
		}, nil
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
		var err error
		t, err = kfsCore.Db.GetEarliestCrated(ctx, hash)
		if err != nil {
			return err
		}
	}
	err := kfsCore.Db.UpsertDCIMMetadataTime(ctx, hash, t)
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
		var err error
		t, err = kfsCore.Db.GetEarliestCrated(ctx, hash)
		if err != nil {
			return err
		}
	}
	err := kfsCore.Db.UpsertDCIMMetadataTime(ctx, hash, t)
	if err != nil {
		return err
	}
	return nil
}

func InsertExif(ctx context.Context, kfsCore *core.KFS, hash string, fileType *dao.FileType) error {
	if fileType.Type == "image" {
		//if fileType.SubType == matchers.TypeJpeg.MIME.Subtype {
		//	GetJpegExifData(kfsCore, hash)
		//}
		width, height, err := getImageHeightWidthByFfmpeg(kfsCore, hash)
		if err != nil {
			return err
		}
		hw := dao.HeightWidth{
			Width:  width,
			Height: height,
		}
		e, err := GetExifData(kfsCore, hash)
		if err != nil {
			err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
			// TODO: what if exist
			if err != nil {
				return err
			}
			_, err = kfsCore.Db.InsertNullExif(ctx, hash)
			// TODO: what if exist
			if err != nil {
				return err
			}
			return nil
		}
		if e.Orientation > 4 {
			hw = dao.HeightWidth{
				Width:  hw.Height,
				Height: hw.Width,
			}
		}
		err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
		// TODO: what if exist
		if err != nil {
			return err
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
		width, height, err := GetVideoHeightWidthByFfmpeg(kfsCore, hash)
		// https://stackoverflow.com/questions/9412384/m4a-mp4-file-format-whats-the-difference-or-are-they-the-same
		if fileType.SubType == "mp4" && err == NO_VALUE_FOR_KEY {
			fileType.Type = "audio"
			return nil
		}
		if err != nil {
			return err
		}
		hw := dao.HeightWidth{
			Width:  width,
			Height: height,
		}
		m, err := GetVideoMetadata(kfsCore, hash)
		if err != nil {
			_, err = kfsCore.Db.InsertNullVideoMetadata(ctx, hash)
			// TODO: what if exist
			if err != nil {
				return err
			}
			return nil
		}
		err = kfsCore.Db.InsertHeightWidth(ctx, hash, hw)
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

var NO_VALUE_FOR_KEY = errors.New("no value for key")

func getValueFor(s string, key string) (string, error) {
	reg, err := regexp.Compile("\n" + key + "=([^\r\n]+)\r?\n")
	if err != nil {
		return "", err
	}
	subMatch := reg.FindStringSubmatch(s)
	if subMatch == nil {
		return "", NO_VALUE_FOR_KEY
	}
	return subMatch[1], nil
}

func getUint64ValueFor(s string, key string) (uint64, error) {
	v, err := getValueFor(s, key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func getInt64ValueFor(s string, key string) (int64, error) {
	v, err := getValueFor(s, key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func getImageHeightWidthByFfmpeg(kfsCore *core.KFS, hash string) (width uint64, height uint64, err error) {
	s, err := ffmpeg_go.Probe(kfsCore.S.GetFilePath(hash), ffmpeg_go.KwArgs{"v": "error", "select_streams": "v", "show_entries": "stream=width,height", "of": "default=noprint_wrappers=1"})
	if err != nil {
		if strings.Contains(err.Error(), "moov atom not found") {
			// https://forum.doom9.org/showthread.php?t=185248
			var rc dao.SizedReadCloser
			rc, err = kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return
			}
			defer rc.Close()
			var c image.Config
			c, err = goheif.DecodeConfig(rc)
			if err != nil {
				return
			}
			height = uint64(c.Height)
			width = uint64(c.Width)
			return
		}
		return
	}
	width, err = getUint64ValueFor(s, "width")
	if err != nil {
		return
	}
	height, err = getUint64ValueFor(s, "height")
	if err != nil {
		return
	}
	rotation, err := getInt64ValueFor(s, "rotation")
	if err != nil {
		return width, height, nil
	}
	// TODO: other rotation
	if rotation == -90 {
		width, height = height, width
	}
	return
}

func GetVideoHeightWidthByFfmpeg(kfsCore *core.KFS, hash string) (width uint64, height uint64, err error) {
	s, err := ffmpeg_go.Probe(kfsCore.S.GetFilePath(hash), ffmpeg_go.KwArgs{"v": "error", "select_streams": "v", "show_entries": "stream=width,height", "of": "default=noprint_wrappers=1"})
	if err != nil {
		return
	}
	width, err = getUint64ValueFor(s, "width")
	if err != nil {
		return
	}
	height, err = getUint64ValueFor(s, "height")
	if err != nil {
		return
	}
	rotation, err := getInt64ValueFor(s, "rotation")
	if err != nil {
		return width, height, nil
	}
	// TODO: other rotation
	if rotation == -90 {
		width, height = height, width
	}
	return
}

func GetVideoMetadata(kfsCore *core.KFS, hash string) (m dao.VideoMetadata, err error) {
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
