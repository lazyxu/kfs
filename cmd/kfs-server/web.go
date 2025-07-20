package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lazyxu/kfs/cmd/kfs-server/task/livp"
	"github.com/lazyxu/kfs/db/sqlc"
	"github.com/lazyxu/kfs/rpc/server"

	"go.uber.org/zap"

	"github.com/lazyxu/kfs/cmd/kfs-server/task/baidu_photo"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/list"
	"github.com/lazyxu/kfs/cmd/kfs-server/task/metadata"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/h2non/filetype/matchers"
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

	e.POST("/api/v1/createFolder", apiCreateFolder)
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
	filePath := c.QueryParams()["filePath[]"]
	var userID int64 = 1
	rc, _, err := kfsCore.OpenFile2(c.Request().Context(), userID, filePath, -1)
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
			img = orientation(img, int64(m.Exif.Orientation))
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
			img = orientation(img, int64(m.Exif.Orientation))
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

func orientation(img image.Image, o int64) image.Image {
	switch o {
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
	idStr := c.QueryParam("id")
	id, err1 := strconv.ParseInt(idStr, 10, 64)
	if err1 != nil {
		return err1
	}
	sizeStr := c.QueryParam("size")
	size, err1 := strconv.ParseInt(sizeStr, 10, 64)
	if err1 != nil {
		return err1
	}
	md5 := c.QueryParam("md5")
	sha256 := c.QueryParam("sha256")
	hwStr := c.QueryParam("hw")
	hw, err1 := strconv.Atoi(hwStr)
	if err1 != nil {
		return err1
	}
	if hw != 64 && hw != 128 && hw != 256 {
		return errors.New("invalid hw, expected 64, 128 or 256")
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
		filename = sizeStr + "_" + md5 + "_" + sha256 + "@" + hwStr + "x" + hwStr
	} else {
		filename = sizeStr + "_" + md5 + "_" + sha256 + "@" + hwStr
	}
	m, err1 := kfsCore.Q.GetMetadata(c.Request().Context(), id)
	if err1 != nil {
		return err1
	}
	thumbnailFilePath := filepath.Join(kfsCore.ThumbnailDir(), filename+".jpg")
	f, err1 := os.Open(thumbnailFilePath)
	if os.IsNotExist(err1) {
		select {
		case thumbnailTaskChan <- struct{}{}:
			println("generate thumbnail for", filename, m.SubType)
		case <-c.Request().Context().Done():
			println("deadlineExceeded generating thumbnail for", filename, m.SubType)
			return context.DeadlineExceeded
		case <-c.Request().Cancel:
			println("canceled generating thumbnail for", filename, m.SubType)
			return context.Canceled
		}
		defer func() {
			<-thumbnailTaskChan
		}()

		if m.Extension == "zip" {
			_, heicHash, err := server.UnzipIfLivp(c.Request().Context(), kfsCore, "")
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
			if m.Orientation.Valid {
				img = orientation(img, m.Orientation.Int64)
			}
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, hw)
			if err != nil {
				return err
			}
		} else if m.Extension == matchers.TypeHeif.Extension {
			rc, err := kfsCore.S2.ReadWithSize(size, md5, sha256)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			if m.Orientation.Valid {
				img = orientation(img, m.Orientation.Int64)
			}
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, hw)
			if err != nil {
				return err
			}
		} else if m.Type == "image" {
			rc, err := kfsCore.S2.ReadWithSize(size, md5, sha256)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := imaging.Decode(rc)
			if err != nil {
				return err
			}
			if m.Orientation.Valid {
				img = orientation(img, m.Orientation.Int64)
			}
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, hw)
			if err != nil {
				return err
			}
		} else if m.Type == "video" {
			originFilePath := filepath.Join(kfsCore.ThumbnailDir(), filename+".origin.jpg")
			src := kfsCore.S2.GetFilePath(size, md5, sha256)
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
			err = generateThumbnail(img, thumbnailFilePath, cutSquare, hw)
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
	dirPathStr := c.QueryParam("dirPath")
	var dirPath []string
	if dirPathStr != "" {
		err := json.Unmarshal([]byte(dirPathStr), &dirPath)
		if err != nil && err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": err.Error()})
		}
	}
	md5 := c.QueryParam("md5")
	sha256 := c.QueryParam("sha256")
	logicFileExist := false
	fileExist, err := kfsCore.Q.GetFileByHash(c.Request().Context(), sqlc.GetFileByHashParams{
		Md5:    md5,
		Sha256: sha256,
	})
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"fileExist": fileExist, "logicFileExist": logicFileExist, "msg": err.Error()})
	}
	var userID int64 = 1
	name := c.QueryParam("name")
	_, files, err := kfsCore.Q.ListDir(c.Request().Context(), userID, dirPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"fileExist": fileExist, "logicFileExist": logicFileExist, "msg": err.Error()})
	}
	for _, f := range files {
		if f.Name == name {
			if f.FileID.Int64 == f.ID {
				logicFileExist = true
				break
			}
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"fileExist": fileExist, "logicFileExist": logicFileExist})
}

func apiCreateFolder(c echo.Context) error {
	type request struct {
		DirPath    []string `json:"dirPath"`
		FolderName string   `json:"folderName"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "无效请求参数"})
	}

	// 调用核心服务创建文件夹
	t := time.Now().UnixNano()
	var userID int64 = 1
	_, err := kfsCore.Q.SetLogicFileByPath(c.Request().Context(), userID, req.DirPath, req.FolderName,
		sql.NullInt64{Valid: false}, sqlc.DefaultDirPerm, 0, t, t, t, t, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

// 上传分块
func apiUploadFileChunk(c echo.Context) error {
	filePathStr := c.FormValue("filePath")
	md5 := c.FormValue("md5")
	sha256 := c.FormValue("sha256")
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
	tmpDir := path.Join(os.TempDir(), "kfs_upload_chunks", md5+"_"+sha256)
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

var defaultFilePerm int64 = 0o644

// 合并分块
func apiMergeFileChunks(c echo.Context) error {
	var req struct {
		DirPath     []string `json:"dirPath"`
		Md5         string   `json:"md5"`
		Sha256      string   `json:"sha256"`
		TotalChunks int      `json:"totalChunks"`
		FileMeta    struct {
			Name         string `json:"name"`
			Size         int64  `json:"size"`
			LastModified int64  `json:"lastModified"`
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
	tmpDir := path.Join(os.TempDir(), "kfs_upload_chunks", req.Md5+"_"+req.Sha256)
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
	log.Printf("merged file: %s", mergedPath)
	stat, err := merged.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "stat merged file failed: "+err.Error())
	}
	size := stat.Size()
	if req.FileMeta.Size != size {
		return c.String(http.StatusInternalServerError, "size is not same: expected "+strconv.Itoa(int(req.FileMeta.Size))+" got "+strconv.Itoa(int(size)))
	}
	ctx := c.Request().Context()

	_, err = kfsCore.S2.Add(req.Md5, req.Sha256, merged)
	if err != nil {
		return c.String(http.StatusInternalServerError, "write to storage failed: "+err.Error())
	}

	// 入库
	res, err := kfsCore.Q.UpsertFile(ctx, sqlc.UpsertFileParams{
		Size:   size,
		Md5:    req.Md5,
		Sha256: req.Sha256,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "insert file failed: "+err.Error())
	}
	fileID, err := res.LastInsertId()
	if err != nil {
		return c.String(http.StatusInternalServerError, "get last insert id failed: "+err.Error())
	}
	var userID int64 = 1
	_, err = kfsCore.Q.SetLogicFileByPath(ctx, userID, req.DirPath, req.FileMeta.Name,
		sql.NullInt64{Valid: true, Int64: fileID}, defaultFilePerm, size,
		req.FileMeta.LastModified, req.FileMeta.LastModified, req.FileMeta.LastModified, req.FileMeta.LastModified, false)
	if err != nil {
		return c.String(http.StatusInternalServerError, "upsert driver file failed: "+err.Error())
	}
	files, _ := os.ReadDir(tmpDir)
	for _, f := range files {
		os.Remove(path.Join(tmpDir, f.Name()))
	}
	// os.RemoveAll(tmpDir)
	return c.String(http.StatusOK, "ok")
}
