package baidu_photo

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/server"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
	kfsCore      *core.KFS
	driverId     uint64
	AccessToken  string
	RefreshToken string
	AppKey       string
	SecretKey    string
	mutex        sync.Locker
	taskInfo     TaskInfo
}

// 目前一刻相册导入相关限制如下：
// 1)大小限制：图片大于50M和视频大于100M的，不支持导入；
// 2)时间限制：图片缺少拍摄时间信息的不支持导入；
// 3)位置限制：仅支持【来自：手机型号】文件夹（如果用户需要导入其他文件夹，可以在网盘里把文件夹移至该里）。
var client = resty.New().
	SetHeader("user-agent", UserAgent).
	SetRetryCount(3).
	SetTimeout(100 * time.Second). // (100M) / (1M/s) = 100s
	SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

func InsertDriverBaiduPhoto(ctx context.Context, kfsCore *core.KFS, name, description, code string) (bool, error) {
	accessToken, refreshToken, err := authByCode(ctx, client, AppKey, SecretKey, code)
	if err != nil {
		return false, err
	}
	exist, err := kfsCore.Db.InsertDriverBaiduPhoto(ctx, name, description, accessToken, refreshToken)
	if err != nil {
		return false, err
	}
	//d := &DriverBaiduPhoto{
	//	kfsCore:      kfsCore,
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//	AppKey:       AppKey,
	//	SecretKey:    SecretKey,
	//}
	//d.test(ctx, name)
	return exist, nil
}

func (d *DriverBaiduPhoto) Get(ctx context.Context, furl string, callback func(req *resty.Request), resp interface{}) error {
	return d.Request(ctx, furl, http.MethodGet, callback, resp, "")
}

func (d *DriverBaiduPhoto) WGet(ctx context.Context, furl string, output string) error {
	return d.Request(ctx, furl, http.MethodGet, nil, nil, output)
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
	req.SetContext(ctx).SetQueryParam("access_token", d.AccessToken)
	if output != "" {
		req.SetOutput(output)
	}
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	res, err := req.Execute(method, furl)
	if err != nil {
		return err
	}

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
	case -6:
		d.AccessToken, d.RefreshToken, err = refreshToken(ctx, client, d.AppKey, d.SecretKey, d.RefreshToken)
		if err != nil {
			return err
		}
		// TODO: save accessToken and refreshToken to db.
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

func (d *DriverBaiduPhoto) Download(ctx context.Context, file File, hash string) error {
	dirPath, name, err := formatPath(file.Path)
	if err != nil {
		return err
	}
	size := uint64(file.Size)
	if hash == "" {
		fmt.Printf("download %s\n", file.Path)
		downloadPath, err := d.GetDownloadPath(ctx, file.Fsid)
		if err != nil {
			return err
		}
		tempFilePath, err := d.DownloadFile(ctx, name, downloadPath)
		if err != nil {
			return err
		}
		f, err := os.Open(tempFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		hash, err = getHash(f)
		if err != nil {
			return err
		}
		_, err = d.kfsCore.S.Write(hash, func(w io.Writer, hasher io.Writer) (e error) {
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
	} else {
		fmt.Printf("exist %s with sha256 %s\n", file.Path, hash)
	}
	err = d.kfsCore.Db.InsertFile(ctx, hash, size)
	if err != nil {
		return err
	}
	err = d.kfsCore.Db.InsertFileMd5(ctx, hash, file.Md5)
	if err != nil {
		return err
	}
	err = d.kfsCore.Db.UpsertDriverFile(ctx, dao.DriverFile{
		DriverId:   d.driverId,
		DirPath:    dirPath,
		Name:       name,
		Version:    0,
		Hash:       hash,
		Mode:       0o777,
		Size:       size,
		CreateTime: uint64(file.ShootTime),
		ModifyTime: uint64(file.ShootTime),
		ChangeTime: uint64(file.ShootTime),
		AccessTime: uint64(file.ShootTime),
	})
	if err != nil {
		return err
	}
	err = server.UpsertLivePhoto(ctx, d.kfsCore, hash, d.driverId, dirPath, name)
	if err != nil {
		return err
	}
	return nil
}

func (d *DriverBaiduPhoto) GetDownloadPath(ctx context.Context, fsid int64) (string, error) {
	var downloadUrl struct {
		Dlink string `json:"dlink"`
	}
	err := d.Get(ctx, FILE_API_URL_V2+"/download", func(r *resty.Request) {
		r.SetQueryParams(map[string]string{
			"fsid": fmt.Sprint(fsid),
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
		List []File `json:"list"`
	}

	File struct {
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

func (d *DriverBaiduPhoto) GetAllFile(ctx context.Context, cb func([]File) bool) error {
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
