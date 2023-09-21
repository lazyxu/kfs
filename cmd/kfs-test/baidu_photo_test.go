package kfs_test

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"testing"
	"time"
)

type BaiduPhoto struct {
	AccessToken  string
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewBaiduPhoto(refreshToken string) *BaiduPhoto {
	return &BaiduPhoto{
		RefreshToken: refreshToken,
		ClientID:     "iYCeC9g08h5vuP9UqvPHKKSVrKFXGa1v",
		ClientSecret: "jXiFMOPVPCWlO2M5CwWQzffpNPaGTRBG",
	}
}

var (
	NoRedirectClient *resty.Client
	RestyClient      *resty.Client
	HttpClient       *http.Client
)

var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
var DefaultTimeout = time.Second * 30

func InitClient() {
	//NoRedirectClient = resty.New().SetRedirectPolicy(
	//	resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//	}),
	//).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: conf.Conf.TlsInsecureSkipVerify})
	//NoRedirectClient.SetHeader("user-agent", UserAgent)

	RestyClient = NewRestyClient()
	//HttpClient = NewHttpClient()
}

var TlsInsecureSkipVerify = true

func NewRestyClient() *resty.Client {
	client := resty.New().
		SetHeader("user-agent", UserAgent).
		SetRetryCount(3).
		SetTimeout(DefaultTimeout).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: TlsInsecureSkipVerify})
	return client
}

type TokenErrResp struct {
	ErrorDescription string `json:"error_description"`
	ErrorMsg         string `json:"error"`
}

var EmptyToken = errors.New("empty token")

func (e *TokenErrResp) Error() string {
	return fmt.Sprint(e.ErrorMsg, " : ", e.ErrorDescription)
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (d *BaiduPhoto) refreshToken() error {
	u := "https://openapi.baidu.com/oauth/2.0/token"
	var resp TokenResp
	var e TokenErrResp
	_, err := RestyClient.R().SetResult(&resp).SetError(&e).SetQueryParams(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": d.RefreshToken,
		"client_id":     d.ClientID,
		"client_secret": d.ClientSecret,
	}).Get(u)
	if err != nil {
		return err
	}
	if e.ErrorMsg != "" {
		return &e
	}
	if resp.RefreshToken == "" {
		return EmptyToken
	}
	d.AccessToken, d.RefreshToken = resp.AccessToken, resp.RefreshToken
	//op.MustSaveDriverStorage(d)
	return nil
}

func (d *BaiduPhoto) Init(ctx context.Context) error {
	if err := d.refreshToken(); err != nil {
		return err
	}
	return nil
}

type Page struct {
	HasMore int    `json:"has_more"`
	Cursor  string `json:"cursor"`
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

const (
	API_URL         = "https://photo.baidu.com/youai"
	USER_API_URL    = API_URL + "/user/v1"
	ALBUM_API_URL   = API_URL + "/album/v1"
	FILE_API_URL_V1 = API_URL + "/file/v1"
	FILE_API_URL_V2 = API_URL + "/file/v2"
)

func (d *BaiduPhoto) Request(furl string, method string, callback func(req *resty.Request), resp interface{}) (*resty.Response, error) {
	req := RestyClient.R().
		SetQueryParam("access_token", d.AccessToken)
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
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
		if err = d.refreshToken(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("errno: %d, refer to https://photo.baidu.com/union/doc", erron)
	}
	return res, nil
}

func (d *BaiduPhoto) Get(furl string, callback func(req *resty.Request), resp interface{}) (*resty.Response, error) {
	return d.Request(furl, http.MethodGet, callback, resp)
}

func (p Page) HasNextPage() bool {
	return p.HasMore == 1
}

// 获取所有文件
func (d *BaiduPhoto) GetAllFile(ctx context.Context) (files []File, err error) {
	var cursor string
	for {
		var resp FileListResp
		_, err = d.Get(FILE_API_URL_V1+"/list", func(r *resty.Request) {
			r.SetContext(ctx)
			r.SetQueryParams(map[string]string{
				"need_thumbnail":     "1",
				"need_filter_hidden": "0",
				"cursor":             cursor,
			})
		}, &resp)
		if err != nil {
			return
		}

		files = append(files, resp.List...)
		if !resp.HasNextPage() {
			return
		}
		cursor = resp.Cursor
	}
}

func TestBaiduPhoto(t *testing.T) {
	// GOPROXY=https://goproxy.cn,direct
	// GOSUMDB=off
	//cmd.Init()

	InitClient()
	// https://alist.nn.ci/zh/guide/drivers/baidu.photo.html
	s := NewBaiduPhoto("122.e2c3359e6741dd988e2989889b4aa30e.Y_aks9WnQqVfUDHEGoFvmXtbz4bISsDxCa731nS.kQHRbA")
	err := s.Init(context.TODO())
	if err != nil {
		t.Error(err)
	}
	files, err := s.GetAllFile(context.TODO())
	if err != nil {
		t.Error(err)
	}
	var size int64
	var cnt int64
	for _, file := range files {
		cnt++
		size += file.Size
	}
	fmt.Printf("size: %s\n", humanize.IBytes(uint64(size)))
	fmt.Printf("cize: %d\n", cnt)

	//var baiduPhoto baidu_photo.BaiduPhoto
	//baiduPhoto.ClientID = "iYCeC9g08h5vuP9UqvPHKKSVrKFXGa1v"
	//baiduPhoto.ClientSecret = "jXiFMOPVPCWlO2M5CwWQzffpNPaGTRBG"
	//baiduPhoto.RefreshToken = "122.238075bc689e8f77bc5388db7991737c.YGu622hbpSoEQh1l4eZx_h87G1BCbqZp60BPXHQ.1tG0sg"
	//baiduPhoto.Init(context.TODO())
	// baiduPhoto.List(context.TODO(), nil, model.ListArgs{})
}
