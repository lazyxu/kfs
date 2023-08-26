package main

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/h2non/filetype/matchers"
	"github.com/jdeng/goheif"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/server"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	e.GET("/api/v1/drivers", apiDrivers)
	e.POST("/api/v1/drivers", apiNewDriver)
	e.DELETE("/api/v1/drivers", apiDeleteDriver)
	e.GET("/api/v1/list", apiList)
	e.GET("/api/v1/openFile", apiOpenFile)
	e.GET("/api/v1/downloadFile", apiDownloadFile)
	e.GET("/api/v1/image", apiImage)

	e.GET("/thumbnail", apiThumbnail)
	e.GET("/api/v1/analysisExif", apiExifStatus)
	e.POST("/api/v1/analysisExif", apiAnalysisExif)
	e.GET("/api/v1/exif", apiListMetadata)
	e.GET("/api/v1/metadata", apiGetMetadata)
	e.GET("/api/v1/diskUsage", apiDiskUsage)

	println("KFS web server listening at:", webPortString)
	// Start server
	e.Logger.Fatal(e.Start(":" + webPortString))
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

func apiDrivers(c echo.Context) error {
	drivers, err := kfsCore.ListDriver(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiNewDriver(c echo.Context) error {
	exist, err := kfsCore.InsertDriver(c.Request().Context(), c.QueryParam("name"), c.QueryParam("description"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, exist)
}

func apiDeleteDriver(c echo.Context) error {
	err := kfsCore.DeleteDriver(c.Request().Context(), c.QueryParam("name"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func apiList(c echo.Context) error {
	driverName := c.QueryParam("driverName")
	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	files, err := kfsCore.ListDriverFile(c.Request().Context(), driverName, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, files)
}

func apiOpenFile(c echo.Context) error {
	driverName := c.QueryParam("driverName")
	filePath := c.QueryParams()["filePath[]"]
	maxContentSizeStr := c.QueryParam("maxContentSize")
	maxContentSize, err := strconv.ParseInt(maxContentSizeStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "maxContentSize should be a number")
	}
	rc, tooLarge, err := kfsCore.OpenFile(c.Request().Context(), driverName, filePath, maxContentSize)
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
	driverName := c.QueryParam("driverName")
	filePath := c.QueryParams()["filePath[]"]
	rc, _, err := kfsCore.OpenFile(c.Request().Context(), driverName, filePath, -1)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}

func apiImage(c echo.Context) error {
	hash := c.QueryParam("hash")
	fileType, err := kfsCore.Db.GetFileType(c.Request().Context(), hash)
	if err != nil {
		return err
	}
	if fileType == server.NewFileType(matchers.TypeHeif) {
		rc, err := kfsCore.S.ReadWithSize(hash)
		if err != nil {
			return err
		}
		defer rc.Close()
		img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusOK
		err = jpeg.Encode(c.Response(), img, &jpeg.Options{Quality: 100})
		if err != nil {
			return err
		}
		return nil
	}
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	defer rc.Close()
	return c.Stream(http.StatusOK, "", rc)
}

func init() {
	err := os.Mkdir("thumbnail", 0o700)
	if os.IsExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
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

func apiThumbnail(c echo.Context) error {
	hash := c.QueryParam("hash")
	sizeStr := c.QueryParam("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return err
	}
	if size != 64 && size != 128 && size != 256 {
		return errors.New("invalid size, expected 64, 128 or 256")
	}
	cutSquareStr := c.QueryParam("cutSquare")
	cutSquare := false
	if cutSquareStr != "" {
		cutSquare, err = strconv.ParseBool(cutSquareStr)
		if err != nil {
			return err
		}
	}
	// TODO: save it to storage.
	var filename string
	if cutSquare {
		filename = hash + "@" + sizeStr + "x" + sizeStr
	} else {
		filename = hash + "@" + sizeStr
	}
	fileType, err := kfsCore.Db.GetFileType(c.Request().Context(), hash)
	if err != nil {
		return err
	}
	thumbnailFilePath := filepath.Join("thumbnail", filename+".jpg")
	f, err := os.Open(thumbnailFilePath)
	if os.IsNotExist(err) {
		println("generate thumbnail for", filename, fileType.SubType)
		if fileType == server.NewFileType(matchers.TypeHeif) {
			var rc dao.SizedReadCloser
			rc, err = kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return err
			}
			defer rc.Close()
			var img image.Image
			img, err = goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, size)
			if err != nil {
				return err
			}
		} else if fileType.Type == "image" {
			var rc dao.SizedReadCloser
			rc, err = kfsCore.S.ReadWithSize(hash)
			if err != nil {
				return err
			}
			defer rc.Close()
			var img image.Image
			img, err = imaging.Decode(rc)
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
		f, err = os.Open(thumbnailFilePath)
	}
	if err != nil {
		return err
	}
	c.Response().Header().Set("Cache-Control", `public, max-age=31536000`)
	defer f.Close()
	return c.Stream(http.StatusOK, "", f)
}
