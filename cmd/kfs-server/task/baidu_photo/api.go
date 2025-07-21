package baidu_photo

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

type TokenErrResp struct {
	ErrorDescription string `json:"error_description"`
	ErrorMsg         string `json:"error"`
}

func (e *TokenErrResp) Error() string {
	return fmt.Sprint(e.ErrorMsg, " : ", e.ErrorDescription)
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DriverBaiduPhoto struct {
	kfsCore    *core.KFS
	id         int64
	BdussBfess string
	Stoken     string
	// AccessToken  string
	// RefreshToken string
	// AppKey       string
	// SecretKey    string
	mutex    sync.Locker
	taskInfo TaskInfo
}

// 目前一刻相册导入相关限制如下：
// 1)大小限制：图片大于50M和视频大于100M的，不支持导入；
// 2)时间限制：图片缺少拍摄时间信息的不支持导入；
// 3)位置限制：仅支持【来自：手机型号】文件夹（如果用户需要导入其他文件夹，可以在网盘里把文件夹移至该里）。

var cookieString string = "BDUSS_BFESS=W1pVXZQSE9KZjhsc2tBS1FjWVFiTFpGaUsyY1h2ZE5mQ1JkZ2RTQzZjNEttNkpvSUFBQUFBJCQAAAAAAAAAAAEAAACamFMptqu3vdDHs71vAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoOe2gKDntoL; STOKEN=8985988b8dff1bb0a9a1e7c7e72546395d5f27a170dc968526034ff086472a02; "

var client = resty.New().
	SetHeader("user-agent", UserAgent).
	SetRetryCount(3).SetRetryWaitTime(3 * time.Second).
	// SetHeader("Cookie", cookieString).
	SetTimeout(100 * time.Second). // (100M) / (1M/s) = 100s
	SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

func InsertDriverBaiduPhoto(ctx context.Context, kfsCore *core.KFS, dirPath []string, name, bdussBfess, stoken, cron string, enabled, syncAtStart bool) error {
	t := time.Now()
	var userID int64 = 1
	id, err := kfsCore.Q.SetLogicFileByPath(ctx, userID, dirPath, name,
		sql.NullInt64{Valid: false}, sqlc.DefaultDirMode, 0, t, t, t, t, false)
	if err != nil {
		return err
	}
	return kfsCore.Q.UpsertSyncFromBaiduPhoto(ctx, sqlc.UpsertSyncFromBaiduPhotoParams{
		ID:          id,
		BdussBfess:  bdussBfess,
		Stoken:      stoken,
		Enabled:     enabled,
		Cron:        cron,
		SyncAtStart: syncAtStart,
	})
}

func (d *DriverBaiduPhoto) Get(ctx context.Context, furl string, callback func(req *resty.Request), resp interface{}) error {
	return d.Request(ctx, furl, http.MethodGet, callback, resp, "")
}

func (d *DriverBaiduPhoto) WGet(ctx context.Context, furl string, output string) error {
	return d.Request(ctx, furl, http.MethodGet, nil, nil, output)
}

func parseCookies(cookieStr string) []*http.Cookie {
	cookies := []*http.Cookie{}
	pairs := strings.Split(cookieStr, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		println("cookie", kv[0], kv[1])
		cookies = append(cookies, &http.Cookie{
			Name:  strings.TrimSpace(kv[0]),
			Value: strings.TrimSpace(kv[1]),
		})
	}
	return cookies
}

func (d *DriverBaiduPhoto) Request(ctx context.Context, furl string, method string, callback func(req *resty.Request), resp interface{}, output string) (err error) {
	req := client.R()
	t := time.Now()
	defer func() {
		if err != nil {
			fmt.Printf("%+v %+v\n", err, time.Now().Sub(t).Seconds())
		} else {
			fmt.Printf("%+v\n", time.Now().Sub(t).Seconds())
		}
	}()
	// req.SetContext(ctx).SetQueryParam("access_token", d.AccessToken)
	req.SetCookie(&http.Cookie{
		Name:  "BDUSS_BFESS",
		Value: d.BdussBfess,
	}).SetCookie(&http.Cookie{
		Name:  "STOKEN",
		Value: d.Stoken,
	})

	if output != "" {
		req.SetOutput(output)
	}
	if callback != nil {
		callback(req)
	}

	// req.SetContext(ctx).SetQueryParam("accessToken", d.AccessToken)
	fmt.Printf("Request %s %s %v\n", method, furl, req.QueryParam)
	if resp != nil {
		req.SetResult(resp)
	}
	res, err := req.Execute(method, furl)
	if err != nil {
		return err
	}
	fmt.Printf("res %v\n", res)
	fmt.Printf("resp %v\n", resp)
	erron := json.ConfigCompatibleWithStandardLibrary.Get(res.Body(), "errno").ToInt()
	switch erron {
	case 0:
		break
	case 50805:
		return fmt.Errorf("you have joined album")
	case 50820:
		return fmt.Errorf("no shared albums found")
	case 50100:
		return fmt.Errorf("illegal title, only supports 50 characters")
	// case -6:
	// 	fmt.Printf("d %v, err: %v\n", d, err)
	// 	d.AccessToken, d.RefreshToken, err = refreshToken(ctx, client, d.AppKey, d.SecretKey, d.RefreshToken)
	// 	fmt.Printf("d %v, err: %v\n", d, err)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = d.kfsCore.Q.UpdateSyncFromBaiduPhoto(ctx, sqlc.UpdateSyncFromBaiduPhotoParams{
	// 		ID:           d.id,
	// 		AccessToken:  d.AccessToken,
	// 		RefreshToken: d.RefreshToken,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// Do not need to goto execute since we have set SetRetryCount to 3.
	default:
		return fmt.Errorf("errno: %d, refer to https://photo.baidu.com/union/doc", erron)
	}
	return nil
}

func getHash(f *os.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func formatPath(filePath string) ([]string, string, error) {
	pathList := strings.Split(filePath, "/")
	newPathList := []string{}
	for _, path := range pathList {
		if path != "" {
			newPathList = append(newPathList, path)
		}
	}
	if len(newPathList) == 0 {
		return newPathList, "", nil
	}
	return newPathList[0 : len(newPathList)-1], newPathList[len(newPathList)-1], nil
}

func (d *DriverBaiduPhoto) Download(ctx context.Context, bdFile BaiduFile, file *sqlc.File) error {
	_, name, err := formatPath(bdFile.Path)
	if err != nil {
		return err
	}
	var sha256 string
	if file == nil {
		fmt.Printf("download %d %s\n", bdFile.Fsid, bdFile.Path)
		downloadPath, err := d.GetDownloadPath(ctx, bdFile.Fsid)
		if err != nil {
			return err
		}
		fmt.Printf("GetDownloadPath %s\n", downloadPath)
		tempFilePath, err := d.DownloadFile(ctx, name, downloadPath)
		if err != nil {
			return err
		}
		f, err := os.Open(tempFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		sha256, err = getHash(f)
		if err != nil {
			return err
		}
		_, err = d.kfsCore.S2.Add(bdFile.Md5, sha256, f)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("exist %s with file %v\n", bdFile.Path, file)
		sha256 = file.Sha256
	}
	res, err := d.kfsCore.Q.UpsertFile(ctx, sqlc.UpsertFileParams{
		Size:   bdFile.Size,
		Md5:    bdFile.Md5,
		Sha256: sha256,
	})
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	createTime := time.Now()
	var userID int64 = 1
	shootTime := time.Unix(bdFile.ShootTime, 0)
	dirPath := []string{shootTime.Format("2006-01")}
	_, err = d.kfsCore.Q.SetLogicFileByParentIdPath(ctx, userID, sqlc.NInt64(d.id), dirPath, name,
		sql.NullInt64{Valid: true, Int64: id}, sqlc.DefaultFileMode, bdFile.Size,
		createTime, shootTime, shootTime, shootTime, true)
	// _, err = d.kfsCore.Q.SetLogicFileByPath(ctx, userID, dirPath, name,
	// 	sql.NullInt64{Valid: true, Int64: id}, sqlc.DefaultFileMode, bdFile.Size,
	// 	createTime, shootTime, shootTime, shootTime, true)
	// err = d.kfsCore.Q.SetLogicFileByPath(ctx, dao.DriverFile{
	// 	DriverId:   d.driverId,
	// 	DirPath:    dirPath,
	// 	Name:       name,
	// 	Version:    0,
	// 	Hash:       hash,
	// 	Mode:       0o777,
	// 	Size:       size,
	// 	CreateTime: uint64(bdFile.ShootTime),
	// 	ModifyTime: uint64(bdFile.ShootTime),
	// 	ChangeTime: uint64(bdFile.ShootTime),
	// 	AccessTime: uint64(bdFile.ShootTime),
	// }, true)
	if err != nil {
		return err
	}
	// err = server.PluginUnzipIfLivp(ctx, d.kfsCore, hash, name)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (d *DriverBaiduPhoto) GetDownloadPath(ctx context.Context, fsid int64) (string, error) {
	var downloadUrl struct {
		Dlink string `json:"dlink"`
	}
	err := d.Get(ctx, FILE_API_URL_V2+"/download", func(r *resty.Request) {
		r.SetQueryParams(map[string]string{
			"fsid":       fmt.Sprint(fsid),
			"clienttype": "70",
		})
	}, &downloadUrl)
	if err != nil {
		return "", err
	}
	return downloadUrl.Dlink, nil
}

func (d *DriverBaiduPhoto) DownloadFile(ctx context.Context, name string, downloadPath string) (string, error) {
	tempDirPath, err := os.MkdirTemp("", "fsid")
	if err != nil {
		return "", err
	}
	tempFilePath := filepath.Join(tempDirPath, name)
	err = d.WGet(ctx, downloadPath, tempFilePath)
	if err != nil {
		return "", err
	}
	return tempFilePath, nil
}

type Page struct {
	HasMore int    `json:"has_more"`
	Cursor  string `json:"cursor"`
}

func (p Page) HasNextPage() bool {
	return p.HasMore == 1
}

type (
	FileListResp struct {
		Page
		List []BaiduFile `json:"list"`
	}

	BaiduFile struct {
		Fsid      int64  `json:"fsid"` // 文件ID
		Path      string `json:"path"` // 文件路径
		Size      int64  `json:"size"`
		Ctime     int64  `json:"ctime"` // 创建时间 s
		Mtime     int64  `json:"mtime"` // 修改时间 s
		Md5       string `json:"md5"`
		ShootTime int64  `json:"shoot_time"`

		//Thumburl []string `json:"thumburl"`
		//parseTime *time.Time
	}
)

func (d *DriverBaiduPhoto) GetAllFile(ctx context.Context, cb func([]BaiduFile) bool) error {
	var cursor string
	for {
		var resp FileListResp
		err := d.Get(ctx, FILE_API_URL_V1+"/list", func(r *resty.Request) {
			r.SetQueryParams(map[string]string{
				"need_thumbnail":     "0",
				"need_filter_hidden": "0",
				//"limit":              "200",
				"cursor": cursor,
			})
		}, &resp)
		if err != nil {
			return err
		}

		continues := cb(resp.List)
		if !continues {
			break
		}

		if !resp.HasNextPage() {
			return nil
		}
		cursor = resp.Cursor
	}
	return nil
}
