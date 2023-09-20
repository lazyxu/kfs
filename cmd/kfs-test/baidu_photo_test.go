package kfs_test

import (
	"context"
	"github.com/alist-org/alist/v3/cmd"
	baidu_photo "github.com/alist-org/alist/v3/drivers/baidu_photo"
	"testing"
)

func TestBaiduPhoto(t *testing.T) {
	// GOPROXY=https://goproxy.cn,direct
	// GOSUMDB=off
	cmd.Init()

	var baiduPhoto baidu_photo.BaiduPhoto
	baiduPhoto.ClientID = "iYCeC9g08h5vuP9UqvPHKKSVrKFXGa1v"
	baiduPhoto.ClientSecret = "jXiFMOPVPCWlO2M5CwWQzffpNPaGTRBG"
	baiduPhoto.RefreshToken = "122.238075bc689e8f77bc5388db7991737c.YGu622hbpSoEQh1l4eZx_h87G1BCbqZp60BPXHQ.1tG0sg"
	baiduPhoto.Init(context.TODO())
	// baiduPhoto.List(context.TODO(), nil, model.ListArgs{})
}
