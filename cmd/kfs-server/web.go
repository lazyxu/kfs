package main

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/lazyxu/kfs/cmd/kfs-server/task/livp"
	"github.com/lazyxu/kfs/rpc/server"

	"go.uber.org/zap"

	"github.com/lazyxu/kfs/cmd/kfs-server/task/baidu_photo"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/list"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/metadata"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/h2non/filetype/matchers"
	"github.com/lazyxu/kfs/dao"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func webServer(webPortString string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           `${time_custom} ${remote_ip} ${method} ${status} ${uri} ${latency_human} ${error}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	// TODO
	zapLogger, _ := zap.NewDevelopment()
	e.Use(Recover(zapLogger))
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
	e.GET("/api/v1/updateDriverLocalFile", apiUpdateDriverLocalFile)
	e.POST("/api/v1/drivers", apiNewDriver)
	e.POST("/api/v1/driverBaiduPhotos", apiNewDriverBaiduPhoto)
	e.POST("/api/v1/driverLocalFiles", apiNewDriverLocalFile)
	e.DELETE("/api/v1/drivers", apiDeleteDriver)
	e.GET("/api/v1/drivers/reset", apiResetDriver)
	e.GET("/api/v1/listLocalFileDriver", apiListLocalFileDriver)

	e.GET("/api/v1/list", apiList)
	e.GET("/api/v1/driverFile", apiGetDriverFile)
	e.GET("/api/v1/event/list", func(c echo.Context) error {
		return list.Handle(c, kfsCore)
	})
	e.GET("/api/v1/listDriverFileByHash", apiListDriverFileByHash)
	e.GET("/api/v1/openFile", apiOpenFile)
	e.GET("/api/v1/downloadFile", apiDownloadFile)
	e.GET("/api/v1/download", apiDownload)
	e.GET("/api/v1/image", apiImage)

	e.GET("/api/v1/drivers/dirCalculatedInfo", apiDriversDirCalculatedInfo)

	e.GET("/thumbnail", apiThumbnail)
	e.GET("/api/v1/analysisExif", apiExifStatus)
	e.POST("/api/v1/analysisExif", apiAnalysisExif)
	e.GET("/api/v1/exif", apiListMetadata)
	e.GET("/api/v1/DCIMMetadataTime", apiListMetadataTime)
	e.GET("/api/v1/listDCIMDriver", apiListDCIMDriver)
	e.GET("/api/v1/listDCIMMediaType", apiListDCIMMediaType)
	e.GET("/api/v1/listDCIMLocation", apiListDCIMLocation)
	e.GET("/api/v1/listDCIMSearchType", apiListDCIMSearchType)
	e.GET("/api/v1/listDCIMSearchSuffix", apiListDCIMSearchSuffix)
	e.GET("/api/v1/searchDCIM", apiSearchDCIM)

	e.GET("/api/v1/metadata", apiGetMetadata)
	e.GET("/api/v1/diskUsage", apiDiskUsage)

	e.POST("/api/v1/startMetadataAnalysisTask", apiStartMetadataAnalysisTask)
	e.GET("/api/v1/event/metadataAnalysisTask", func(c echo.Context) error {
		return metadata.ApiEvent(c, kfsCore)
	})
	e.POST("/api/v1/startLivePhotoAnalysisTask", apiStartLivePhotoAnalysisTask)
	e.GET("/api/v1/event/livePhotoAnalysisTask", func(c echo.Context) error {
		return livp.ApiEvent(c, kfsCore)
	})
	e.POST("/api/v1/startBaiduPhotoTask", apiStartBaiduPhotoTask)
	e.GET("/api/v1/event/baiduPhotoTask/:driverId", func(c echo.Context) error {
		return baidu_photo.ApiEvent(c, kfsCore)
	}) // TODO: handle name

	e.POST("/api/v1/uploadFileChunk", apiUploadFileChunk)
	e.POST("/api/v1/mergeFileChunks", apiMergeFileChunks)
	e.GET("/api/v1/checkFileExists", apiCheckFileExists)

	println("KFS web server listening at:", webPortString)
	// Start server
	e.Logger.Fatal(e.Start(":" + webPortString))
}

func apiDriversDirCalculatedInfo(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	n, err := kfsCore.Db.GetDriverDirCalculatedInfo(c.Request().Context(), driverId, filePath)
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
	files, err := kfsCore.Db.ListDriverFile(c.Request().Context(), driverId, filePath)
	if err != nil {
		println(err.Error())
		c.Logger().Error(err)
		return err
	}
	return ok(c, files)
}

func apiGetDriverFile(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, "driverId should be a number")
	}
	filePath := c.QueryParams()["filePath[]"]
	if filePath == nil {
		filePath = []string{}
	}
	files, err := kfsCore.Db.GetDriverFile(c.Request().Context(), driverId, filePath)
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
	if fileType.Extension == "zip" {
		thumbnailFilePath := filepath.Join(kfsCore.TransCodeDir(), hash+".jpg")
		f, err2 := os.Open(thumbnailFilePath)
		if os.IsNotExist(err2) {
			_, heicHash, err := server.UnzipIfLivp(c.Request().Context(), kfsCore, hash)
			if err != nil {
				return err
			}
			rc, err := kfsCore.S.ReadWithSize(heicHash)
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
	}
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
			err = ffmpeg_go.Input(src, ffmpeg_go.KwArgs{"noautorotate": ""}).
				Output(thumbnailFilePath, ffmpeg_go.KwArgs{"qscale": "0", "vcodec": "libx264"}).
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

var thumbnailTaskChan = make(chan struct{}, 1)

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
		select {
		case thumbnailTaskChan <- struct{}{}:
			println("generate thumbnail for", filename, fileType.SubType)
		case <-c.Request().Context().Done():
			println("deadlineExceeded generating thumbnail for", filename, fileType.SubType)
			return context.DeadlineExceeded
		case <-c.Request().Cancel:
			println("canceled generating thumbnail for", filename, fileType.SubType)
			return context.Canceled
		}
		defer func() {
			<-thumbnailTaskChan
		}()

		if fileType.Extension == "zip" {
			_, heicHash, err := server.UnzipIfLivp(c.Request().Context(), kfsCore, hash)
			if err != nil {
				return err
			}
			rc, err := kfsCore.S.ReadWithSize(heicHash)
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
		} else if fileType.Extension == matchers.TypeHeif.Extension {
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
				Output(originFilePath, ffmpeg_go.KwArgs{"frames:v": 1, "loglevel": "error"}).
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
	forceStr := c.QueryParam("force")
	force, _ := strconv.ParseBool(forceStr)
	metadata.StartOrStop(kfsCore, start, force)
	return c.String(http.StatusOK, "")
}

func apiStartLivePhotoAnalysisTask(c echo.Context) error {
	startStr := c.QueryParam("start")
	start, err := strconv.ParseBool(startStr)
	if err != nil {
		return err
	}
	forceStr := c.QueryParam("force")
	force, _ := strconv.ParseBool(forceStr)
	livp.StartOrStop(kfsCore, start, force)
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

// 检查文件是否已存在
func apiCheckFileExists(c echo.Context) error {
	driverIdStr := c.QueryParam("driverId")
	driverId, err := strconv.ParseUint(driverIdStr, 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"exists": false, "msg": "driverId should be a number"})
	}
	filePathStr := c.QueryParam("filePath")
	var filePath []string
	if filePathStr != "" {
		_ = json.Unmarshal([]byte(filePathStr), &filePath)
	}
	hash := c.QueryParam("hash")
	// 检查driver+filePath下是否有hash为该值的文件
	files, err := kfsCore.Db.ListDriverFile(c.Request().Context(), driverId, filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"exists": false, "msg": err.Error()})
	}
	for _, f := range files {
		if f.Hash == hash {
			return c.JSON(http.StatusOK, map[string]interface{}{"exists": true})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"exists": false})
}

// 上传分块
func apiUploadFileChunk(c echo.Context) error {
	filePathStr := c.FormValue("filePath")
	hash := c.FormValue("hash")
	chunkIndex := c.FormValue("chunkIndex")
	file, err := c.FormFile("chunk")
	if err != nil {
		return c.String(http.StatusBadRequest, "missing chunk")
	}
	var filePath []string
	if filePathStr != "" {
		_ = json.Unmarshal([]byte(filePathStr), &filePath)
	}
	// 文件元信息
	fileName := c.FormValue("fileName")
	fileSize := c.FormValue("fileSize")
	fileLastModified := c.FormValue("fileLastModified")
	fileType := c.FormValue("fileType")
	// 可将这些信息写入临时目录的meta.json
	tmpDir := path.Join(os.TempDir(), "kfs_upload_chunks", hash)
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return c.String(http.StatusInternalServerError, "mkdir failed: "+err.Error())
	}
	// 写meta.json（只写一次）
	metaPath := path.Join(tmpDir, "meta.json")
	if fileName != "" && fileSize != "" && fileLastModified != "" {
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			meta := map[string]string{
				"name":         fileName,
				"size":         fileSize,
				"lastModified": fileLastModified,
				"type":         fileType,
			}
			bs, _ := json.Marshal(meta)
			_ = os.WriteFile(metaPath, bs, 0644)
		}
	}
	// ...existing code...
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "open chunk failed: "+err.Error())
	}
	defer src.Close()
	chunkPath := path.Join(tmpDir, chunkIndex)
	dst, err := os.Create(chunkPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "create chunk failed: "+err.Error())
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "write chunk failed: "+err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

// 合并分块
func apiMergeFileChunks(c echo.Context) error {
	var req struct {
		DriverId    uint64   `json:"driverId"`
		DirPath     []string `json:"dirPath"`
		Hash        string   `json:"hash"`
		TotalChunks int      `json:"totalChunks"`
		FileMeta    struct {
			Name         string `json:"name"`
			Size         uint64 `json:"size"`
			LastModified uint64 `json:"lastModified"`
			Type         string `json:"type"`
		} `json:"fileMeta"`
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "read body failed: "+err.Error())
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return c.String(http.StatusBadRequest, "json decode failed: "+err.Error())
	}
	tmpDir := path.Join(os.TempDir(), "kfs_upload_chunks", req.Hash)
	mergedPath := path.Join(tmpDir, "merged_file")
	merged, err := os.Create(mergedPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "create merged file failed: "+err.Error())
	}
	defer merged.Close()
	// 按顺序合并
	for i := 0; i < req.TotalChunks; i++ {
		chunkPath := path.Join(tmpDir, strconv.Itoa(i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "open chunk failed: "+err.Error())
		}
		_, err = io.Copy(merged, chunk)
		chunk.Close()
		if err != nil {
			return c.String(http.StatusInternalServerError, "merge chunk failed: "+err.Error())
		}
	}
	// 读取meta.json（如果有）
	metaPath := path.Join(tmpDir, "meta.json")
	var meta map[string]string
	if bs, err := os.ReadFile(metaPath); err == nil {
		_ = json.Unmarshal(bs, &meta)
	}
	// TODO: 将mergedPath文件保存到最终存储，并入库
	// 这里只做简单演示，实际应调用kfsCore相关接口
	// 清理分块
	merged.Seek(0, io.SeekStart)
	stat, err := merged.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "stat merged file failed: "+err.Error())
	}
	size := uint64(stat.Size())
	ctx := c.Request().Context()

	// 保存文件到存储
	_, err = merged.Seek(0, io.SeekStart)
	if err != nil {
		return c.String(http.StatusInternalServerError, "seek merged file failed: "+err.Error())
	}
	_, err = kfsCore.S.Write(req.Hash, func(w io.Writer, hasher io.Writer) (e error) {
		rr := io.TeeReader(merged, hasher)
		_, err = io.Copy(w, rr)
		return err
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "write to storage failed: "+err.Error())
	}

	// 入库
	err = kfsCore.Db.InsertFile(ctx, req.Hash, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, "insert file failed: "+err.Error())
	}
	err = kfsCore.Db.UpsertDriverFile(ctx, dao.DriverFile{
		DriverId:   req.DriverId,
		DirPath:    req.DirPath,
		Hash:       req.Hash,
		Name:       req.FileMeta.Name,
		Version:    0,
		Mode:       0o644,
		Size:       size,
		CreateTime: req.FileMeta.LastModified,
		ModifyTime: req.FileMeta.LastModified,
		ChangeTime: req.FileMeta.LastModified,
		AccessTime: req.FileMeta.LastModified,
	}, false)
	if err != nil {
		return c.String(http.StatusInternalServerError, "upsert driver file failed: "+err.Error())
	}
	files, _ := ioutil.ReadDir(tmpDir)
	for _, f := range files {
		os.Remove(path.Join(tmpDir, f.Name()))
	}
	os.RemoveAll(tmpDir)
	return c.String(http.StatusOK, "ok")
}
