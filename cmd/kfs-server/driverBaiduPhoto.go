package main

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"time"
)

type DriverBaiduPhoto struct {
	AccessToken  string
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
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

func AuthByCode(code string) (string, string, error) {
	return authByCode(client, AppKey, SecretKey, code)
}

func authByCode(client *resty.Client, appKey string, secretKey string, code string) (string, string, error) {
	u := "https://openapi.baidu.com/oauth/2.0/token"
	var resp TokenResp
	var e TokenErrResp
	_, err := client.R().SetResult(&resp).SetError(&e).SetQueryParams(map[string]string{
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
