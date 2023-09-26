package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/server"
)

type DriverBaiduPhoto struct {
	AccessToken  string
	RefreshToken string
	AppKey       string
	SecretKey    string
}

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
const DefaultTimeout = time.Second * 30

const AppKey = "huREKC2eNTctaBWfh3LdiAYjZ9ARBh5g"
const SecretKey = "eMmhaLDpxzTKX3upCguM0q9yOsmVDP6g"

var client = resty.New().
	SetHeader("user-agent", UserAgent).
	SetRetryCount(3).
	SetTimeout(DefaultTimeout).
	SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

func InsertDriverBaiduPhoto(ctx context.Context, name, description, typ, code string) (bool, error) {
	accessToken, refreshToken, err := authByCode(ctx, client, AppKey, SecretKey, code)
	if err != nil {
		return false, err
	}
	exist, err := kfsCore.Db.InsertDriver(ctx, name, description, typ, accessToken, refreshToken)
	if err != nil {
		return false, err
	}
	d := &DriverBaiduPhoto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AppKey:       AppKey,
		SecretKey:    SecretKey,
	}
	d.test(ctx, name)
	return exist, nil
}

func (d *DriverBaiduPhoto) Get(ctx context.Context, furl string, callback func(req *resty.Request), resp interface{}) (*resty.Response, error) {
	return d.Request(ctx, furl, http.MethodGet, callback, resp)
}

func (d *DriverBaiduPhoto) Request(ctx context.Context, furl string, method string, callback func(req *resty.Request), resp interface{}) (*resty.Response, error) {
	req := client.R()
	req.SetContext(ctx).SetQueryParam("access_token", d.AccessToken)
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	var refreshed bool
execute:
	res, err := req.Execute(method, furl)
	if err != nil {
		return nil, err
	}

	erron := json.ConfigCompatibleWithStandardLibrary.Get(res.Body(), "errno").ToInt()
	switch erron {
	case 0:
		break
	case 50805:
		return nil, fmt.Errorf("you have joined album")
	case 50820:
		return nil, fmt.Errorf("no shared albums found")
	case 50100:
		return nil, fmt.Errorf("illegal title, only supports 50 characters")
	case -6:
		if refreshed {
			return nil, fmt.Errorf("invalid token after refeshed")
		}
		d.AccessToken, d.RefreshToken, err = refreshToken(ctx, client, d.AppKey, d.SecretKey, d.RefreshToken)
		if err != nil {
			return nil, err
		}
		refreshed = true
		// TODO: save accessToken and refreshToken to db.
		goto execute
	default:
		return nil, fmt.Errorf("errno: %d, refer to https://photo.baidu.com/union/doc", erron)
	}
	return res, nil
}

const (
	API_URL         = "https://photo.baidu.com/youai"
	USER_API_URL    = API_URL + "/user/v1"
	ALBUM_API_URL   = API_URL + "/album/v1"
	FILE_API_URL_V1 = API_URL + "/file/v1"
	FILE_API_URL_V2 = API_URL + "/file/v2"
)

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

func (d *DriverBaiduPhoto) Download(ctx context.Context, file File, driverName string) error {
	downloadPath, err := d.GetDownloadPath(ctx, file.Fsid)
	if err != nil {
		return err
	}
	req := client.R()
	tempDirPath, err := os.MkdirTemp("", "fsid")
	if err != nil {
		return err
	}
	dirPath, name, err := formatPath(file.Path)
	if err != nil {
		return err
	}
	tempFilePath := filepath.Join(tempDirPath, name)
	_, err = req.SetContext(ctx).SetOutput(tempFilePath).Get(downloadPath)
	if err != nil {
		return err
	}
	f, err := os.Open(tempFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	size := uint64(info.Size())
	hash, err := getHash(f)
	if err != nil {
		return err
	}
	_, err = kfsCore.S.Write(hash, func(w io.Writer, hasher io.Writer) (e error) {
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		rr := io.TeeReader(f, hasher)
		_, err = io.Copy(w, rr)
		return err
	})
	if err != nil {
		return err
	}
	err = kfsCore.Db.InsertFile(ctx, hash, size)
	if err != nil {
		return err
	}
	err = kfsCore.Db.UpsertDriverFile(context.TODO(), dao.DriverFile{
		DriverName: driverName,
		DirPath:    dirPath,
		Name:       name,
		Version:    0,
		Hash:       hash,
		Mode:       0o777,
		Size:       size,
		CreateTime: uint64(file.Ctime),
		ModifyTime: uint64(file.Mtime),
		ChangeTime: uint64(file.Mtime),
		AccessTime: uint64(file.Mtime),
	})
	if err != nil {
		return err
	}
	err = server.UpsertLivePhoto(kfsCore, hash, driverName, dirPath, name)
	if err != nil {
		return err
	}
	return nil
}

func (d *DriverBaiduPhoto) GetDownloadPath(ctx context.Context, fsid int64) (string, error) {
	var downloadUrl struct {
		Dlink string `json:"dlink"`
	}
	_, err := d.Get(ctx, FILE_API_URL_V2+"/download", func(r *resty.Request) {
		r.SetQueryParams(map[string]string{
			"fsid": fmt.Sprint(fsid),
		})
	}, &downloadUrl)
	if err != nil {
		return "", err
	}
	return downloadUrl.Dlink, nil
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
		List []File `json:"list"`
	}

	File struct {
		Fsid     int64    `json:"fsid"` // 文件ID
		Path     string   `json:"path"` // 文件路径
		Size     int64    `json:"size"`
		Ctime    int64    `json:"ctime"` // 创建时间 s
		Mtime    int64    `json:"mtime"` // 修改时间 s
		Thumburl []string `json:"thumburl"`

		parseTime *time.Time
	}
)

func (d *DriverBaiduPhoto) test(ctx context.Context, driverName string) {
	d.GetAllFile(ctx, func(list []File) bool {
		for i, f := range list {
			fmt.Printf("[%d/%d] downloading %s\n", i, len(list), f.Path)
			d.Download(ctx, f, driverName)
		}
		return true
	})
}

func (d *DriverBaiduPhoto) GetAllFile(ctx context.Context, cb func([]File) bool) error {
	var cursor string
	for {
		var resp FileListResp
		_, err := d.Get(ctx, FILE_API_URL_V1+"/list", func(r *resty.Request) {
			r.SetQueryParams(map[string]string{
				"need_thumbnail":     "1",
				"need_filter_hidden": "0",
				"cursor":             cursor,
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

func authByCode(ctx context.Context, client *resty.Client, appKey string, secretKey string, code string) (string, string, error) {
	u := "https://openapi.baidu.com/oauth/2.0/token"
	var resp TokenResp
	var e TokenErrResp
	req := client.R()
	_, err := req.SetContext(ctx).SetResult(&resp).SetError(&e).SetQueryParams(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     appKey,
		"client_secret": secretKey,
		"redirect_uri":  "oob",
	}).Get(u)
	if err != nil {
		return "", "", err
	}
	if e.ErrorMsg != "" {
		return "", "", &e
	}
	if resp.RefreshToken == "" {
		return "", "", EmptyToken
	}
	return resp.AccessToken, resp.RefreshToken, nil
}

func refreshToken(ctx context.Context, client *resty.Client, appKey string, secretKey string, refreshToken string) (string, string, error) {
	u := "https://openapi.baidu.com/oauth/2.0/token"
	var resp TokenResp
	var e TokenErrResp
	req := client.R()
	_, err := req.SetContext(ctx).SetResult(&resp).SetError(&e).SetQueryParams(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     appKey,
		"client_secret": secretKey,
	}).Get(u)
	if err != nil {
		return "", "", err
	}
	if e.ErrorMsg != "" {
		return "", "", &e
	}
	if resp.RefreshToken == "" {
		return "", "", EmptyToken
	}
	return resp.AccessToken, resp.RefreshToken, nil
}
