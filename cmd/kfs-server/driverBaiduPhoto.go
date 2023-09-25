package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"time"
)

type DriverBaiduPhoto struct {
	AccessToken  string
	RefreshToken string `json:"refresh_token"`
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

func (d *DriverBaiduPhoto) Download(ctx context.Context, fsid int64) error {
	downloadPath, err := d.GetDownloadPath(ctx, fsid)
	if err != nil {
		return err
	}
	kfsCore.S.Write()
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
