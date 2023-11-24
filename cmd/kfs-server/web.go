package main

import (
	"errors"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/baidu_photo"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/metadata"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/h2non/filetype/matchers"
	"github.com/jdeng/goheif"
	"github.com/lazyxu/kfs/dao"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func webServer(webPortString string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	e.StaticFS("/", echo.MustSubFS(build, "build"))

	// Routes
	e.GET("/api/v1/devices", apiDevices)
	e.POST("/api/v1/devices", apiNewDevice)
	e.DELETE("/api/v1/devices", apiDeleteDevice)

	e.GET("/api/v1/drivers", apiDrivers)
	e.GET("/api/v1/getDriverSync", apiGetDriverSync)
	e.GET("/api/v1/getDriverLocalFile", apiGetDriverLocalFile)
	e.GET("/api/v1/updateDriverSync", apiUpdateDriverSync)
	e.POST("/api/v1/drivers", apiNewDriver)
	e.POST("/api/v1/driverBaiduPhotos", apiNewDriverBaiduPhoto)
	e.POST("/api/v1/driverLocalFiles", apiNewDriverLocalFile)
	e.DELETE("/api/v1/drivers", apiDeleteDriver)
	e.GET("/api/v1/drivers/reset", apiResetDriver)
	e.GET("/api/v1/listLocalFileDriver", apiListLocalFileDriver)

	e.GET("/api/v1/list", apiList)
	e.GET("/api/v1/listDriverFileByHash", apiListDriverFileByHash)
	e.GET("/api/v1/openFile", apiOpenFile)
	e.GET("/api/v1/downloadFile", apiDownloadFile)
	e.GET("/api/v1/download", apiDownload)
	e.GET("/api/v1/image", apiImage)

	e.GET("/api/v1/drivers/fileSize", apiDriversFileSize)
	e.GET("/api/v1/drivers/fileCount", apiDriversFileCount)
	e.GET("/api/v1/drivers/dirCount", apiDriversDirCount)

	e.GET("/thumbnail", apiThumbnail)
	e.GET("/api/v1/analysisExif", apiExifStatus)
	e.POST("/api/v1/analysisExif", apiAnalysisExif)
	e.GET("/api/v1/exif", apiListMetadata)
	e.GET("/api/v1/metadata", apiGetMetadata)
	e.GET("/api/v1/diskUsage", apiDiskUsage)

	e.POST("/api/v1/startMetadataAnalysisTask", apiStartMetadataAnalysisTask)
	e.GET("/api/v1/event/metadataAnalysisTask", func(c echo.Context) error {
		return metadata.ApiEvent(c, kfsCore)
	})
	e.POST("/api/v1/startBaiduPhotoTask", apiStartBaiduPhotoTask)
	e.GET("/api/v1/event/baiduPhotoTask/:driverId", func(c echo.Context) error {
		return baidu_photo.ApiEvent(c, kfsCore)
	}) // TODO: handle name

	println("KFS web server listening at:", webPortString)
	// Start server
	e.Logger.Fatal(e.Start(":" + webPortString))
}

func apiDriversFileSize(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	n, err := kfsCore.Db.GetDriverFileSize(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, n)
}

func apiDriversFileCount(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	n, err := kfsCore.Db.GetDriverFileCount(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, n)
}

func apiDriversDirCount(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	n, err := kfsCore.Db.GetDriverDirCount(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, n)
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func ok(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{Data: data})
}

// Handler

func apiList(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	files, err := kfsCore.ListDriverFile(c.Request().Context(), driverId, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, files)
}

func apiOpenFile(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	maxContentSizeStr := c.QueryParam("maxContentSize")
	maxContentSize, err := strconv.ParseInt(maxContentSizeStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "maxContentSize should be a number")
	}
	rc, tooLarge, err := kfsCore.OpenFile(c.Request().Context(), driverId, filePath, maxContentSize)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	if tooLarge {
		c.Response().Header().Set("Kfs-Too-Large", "true")
		return c.String(http.StatusOK, "")
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}

func apiDownloadFile(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "id should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	rc, _, err := kfsCore.OpenFile(c.Request().Context(), driverId, filePath, -1)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}

func apiDownload(c echo.Context) error {
	hash := c.QueryParam("hash")
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
	return c.Stream(http.StatusOK, "", rc)
}

func apiListDriverFileByHash(c echo.Context) error {
	hash := c.QueryParam("hash")
	list, err := kfsCore.Db.ListDriverFileByHash(c.Request().Context(), hash)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, list)
}

func apiImage(c echo.Context) error {
	hash := c.QueryParam("hash")
	m, err1 := kfsCore.Db.GetMetadata(c.Request().Context(), hash)
	if err1 != nil {
		return err1
	}
	fileType := m.FileType
	if fileType.Extension == matchers.TypeHeif.Extension {
		thumbnailFilePath := filepath.Join(kfsCore.TransCodeDir(), hash+".jpg")
		f, err2 := os.Open(thumbnailFilePath)
		if os.IsNotExist(err2) {
			rc, err := kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			img = orientation(img, m.Exif)
			err = imaging.Save(img, thumbnailFilePath)
			if err != nil {
				return err
			}
			f, err2 = os.Open(thumbnailFilePath)
		}
		if err2 != nil {
			return err2
		}
		defer f.Close()
		c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
		return c.Stream(http.StatusOK, "", f)
	} else if fileType.Extension == matchers.TypeMov.Extension {
		src := kfsCore.S.GetFilePath(hash)
		thumbnailFilePath := filepath.Join(kfsCore.TransCodeDir(), hash+".mp4")
		f, err := os.Open(thumbnailFilePath)
		if os.IsNotExist(err) {
			err = ffmpeg_go.Input(src).
				Output(thumbnailFilePath, ffmpeg_go.KwArgs{"qscale": "0"}).
				OverWriteOutput().ErrorToStdOut().Run()
			if err != nil {
				return err
			}
			f, err = os.Open(thumbnailFilePath)
		}
		if err != nil {
			return err
		}
		defer f.Close()
		c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
		return c.Stream(http.StatusOK, "", f)
	}
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
	return c.Stream(http.StatusOK, "", rc)
}

func generateThumbnail(img image.Image, thumbnailFilePath string, cutSquare bool, size int) error {
	x := img.Bounds().Size().X
	y := img.Bounds().Size().Y
	var newImg *image.NRGBA
	if cutSquare {
		newImg = imaging.Thumbnail(img, size, size, imaging.Lanczos)
	} else {
		var xx int
		var yy int
		if x > y {
			xx = size
			yy = int(float64(size) * float64(y) / float64(x))
		} else {
			xx = int(float64(size) * float64(x) / float64(y))
			yy = size
		}
		newImg = imaging.Thumbnail(img, xx, yy, imaging.Lanczos)
	}
	err := imaging.Save(newImg, thumbnailFilePath)
	if err != nil {
		return err
	}
	return nil
}

func orientation(img image.Image, exif *dao.Exif) image.Image {
	if exif == nil {
		return img
	}
	switch exif.Orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.Rotate180(imaging.FlipH(img))
	case 5:
		return imaging.Rotate90(imaging.FlipH(img))
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Rotate270(imaging.FlipH(img))
	case 8:
		return imaging.Rotate90(img)
	}
	return img
}

func apiThumbnail(c echo.Context) error {
	hash := c.QueryParam("hash")
	sizeStr := c.QueryParam("size")
	size, err1 := strconv.Atoi(sizeStr)
	if err1 != nil {
		return err1
	}
	if size != 64 && size != 128 && size != 256 {
		return errors.New("invalid size, expected 64, 128 or 256")
	}
	cutSquareStr := c.QueryParam("cutSquare")
	cutSquare := false
	if cutSquareStr != "" {
		cutSquare, err1 = strconv.ParseBool(cutSquareStr)
		if err1 != nil {
			return err1
		}
	}
	// TODO: save it to storage.
	var filename string
	if cutSquare {
		filename = hash + "@" + sizeStr + "x" + sizeStr
	} else {
		filename = hash + "@" + sizeStr
	}
	m, err1 := kfsCore.Db.GetMetadata(c.Request().Context(), hash)
	if err1 != nil {
		return err1
	}
	fileType := m.FileType
	thumbnailFilePath := filepath.Join(kfsCore.ThumbnailDir(), filename+".jpg")
	f, err1 := os.Open(thumbnailFilePath)
	if os.IsNotExist(err1) {
		println("generate thumbnail for", filename, fileType.SubType)
		if fileType.Extension == matchers.TypeHeif.Extension {
			rc, err := kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			img = orientation(img, m.Exif)
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, size)
			if err != nil {
				return err
			}
		} else if fileType.Type == "image" {
			rc, err := kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := imaging.Decode(rc)
			if err != nil {
				return err
			}
			img = orientation(img, m.Exif)
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, size)
			if err != nil {
				return err
			}
		} else if fileType.Type == "video" {
			originFilePath := filepath.Join(kfsCore.ThumbnailDir(), filename+".origin.jpg")
			src := kfsCore.S.GetFilePath(hash)
			err := ffmpeg_go.Input(src).
				Output(originFilePath, ffmpeg_go.KwArgs{"vframes": 1}).
				OverWriteOutput().ErrorToStdOut().Run()
			if err != nil {
				return err
			}
			f, err := os.Open(originFilePath)
			if err != nil {
				return err
			}
			defer f.Close()
			img, err := imaging.Decode(f)
			if err != nil {
				return err
			}
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, size)
			if err != nil {
				return err
			}
		} else {
			return errors.New("unsupported file type")
		}
		f, err1 = os.Open(thumbnailFilePath)
	}
	if err1 != nil {
		return err1
	}
	c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
	defer f.Close()
	return c.Stream(http.StatusOK, "", f)
}

func apiStartMetadataAnalysisTask(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	metadata.StartOrStop(kfsCore, start)
	return c.String(http.StatusOK, "")
}

func apiStartBaiduPhotoTask(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	ctx := c.Request().Context()
	d, err := baidu_photo.GetOrLoadDriver(ctx, kfsCore, driverId)
	if err != nil {
		return err
	}
	d.StartOrStop(ctx, start)
	return c.String(http.StatusOK, "")
}
