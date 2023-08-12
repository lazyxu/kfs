package main

import (
	"github.com/disintegration/imaging"
	"github.com/lazyxu/kfs/dao"
	"image"
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

	e.GET("/thumbnail", apiThumbnail)

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
	drivers, err := kfsCore.DriverList(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return ok(c, drivers)
}

func apiNewDriver(c echo.Context) error {
	exist, err := kfsCore.NewDriver(c.Request().Context(), c.QueryParam("name"), c.QueryParam("description"))
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
	files, err := kfsCore.ListV2(c.Request().Context(), driverName, filePath)
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
func init() {
	err := os.Mkdir("thumbnail", 0o700)
	if os.IsExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
}

func apiThumbnail(c echo.Context) error {
	hash := c.QueryParam("hash")
	thumbnailFilePath := filepath.Join("thumbnail", hash+".jpg")
	f, err := os.Open(thumbnailFilePath)
	if os.IsNotExist(err) {
		var rc dao.SizedReadCloser
		rc, err = kfsCore.S.ReadWithSize(hash)
		if err != nil {
			println(err.Error())
			c.Logger().Error(err)
			return err
		}
		defer rc.Close()
		var img image.Image
		img, err = imaging.Decode(rc)
		x := img.Bounds().Size().X
		y := img.Bounds().Size().Y
		var xx int
		var yy int
		if x > y {
			xx = 64
			yy = int(64.0 * float64(y) / float64(x))
		} else {
			xx = int(64.0 * float64(x) / float64(y))
			yy = 64
		}
		newImg := imaging.Resize(img, xx, yy, imaging.Lanczos)
		err = imaging.Save(newImg, thumbnailFilePath)
		if err != nil {
			return err
		}
		f, err = os.Open(thumbnailFilePath)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Stream(http.StatusOK, "", f)
}
