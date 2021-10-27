package aliyunpan

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/library-go/logger"
)

type Client struct {
	*aliyunpan.PanClient
	UserInfo *aliyunpan.UserInfo
}

func LoginByRefreshToken(refreshToken string) *Client {
	webToken, _ := aliyunpan.GetAccessTokenFromRefreshToken(refreshToken)
	fmt.Println(webToken)
	if webToken == nil {
		return nil
	}

	panClient := aliyunpan.NewPanClient(*webToken, aliyunpan.AppLoginToken{})

	fmt.Println(" ")
	ui, _ := panClient.GetUserInfo()
	fmt.Println(ui)
	return &Client{
		PanClient: panClient,
		UserInfo:  ui,
	}
}

func (client *Client) Init() {
	err := client.upload("kfs-config.json")
	if err != nil {
		fmt.Println("InitError", err)
	}
}

func (client *Client) upload(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return err
	}
	logger.IsVerbose = true
	x, e := client.CreateUploadFile(&aliyunpan.CreateFileUploadParam{
		Name:         name,
		DriveId:      client.UserInfo.FileDriveId,
		ParentFileId: aliyunpan.DefaultRootParentFileId,
		Size:         fi.Size(),
		ContentHash:  string(hash.Sum(nil)),
	})
	if e != nil {
		return e
	}
	fmt.Println("upload", x)
	return nil
}

func (client *Client) ls() {
	fl, _ := client.FileList(&aliyunpan.FileListParam{
		DriveId:      client.UserInfo.FileDriveId,
		ParentFileId: aliyunpan.DefaultRootParentFileId,
		Limit:        10,
	})
	fmt.Println(fl)
}
