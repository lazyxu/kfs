package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lazyxu/kfs/cmd/client/kfsclient"
)

// {
//   "clientID": "3fb2f545-a11e-409f-ad8e-f3bcc35bfcd0",
//   "theme": "dark",
//   "backendProcess": {
//     "port": "1123",
//     "status": "运行中"
//   },
//   "username": "17161951517",
//   "refreshToken": "96246b97eb994fcaa4e8abb553d502bb",
//   "downloadPath": ""
// }

type Config struct {
	ClientID string
}

func GetConfig() (*Config, error) {
	file, err := ioutil.ReadFile("kfs-config.json")
	if err != nil {
		return nil, err
	}
	data := &Config{}
	err = json.Unmarshal(file, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

var gClient *kfsclient.Client
var once sync.Once

func GetClient() *kfsclient.Client {
	once.Do(func() {
		gClient = kfsclient.New("localhost:9092")
	})
	return gClient
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, this is kfs client!")
	})
	e.GET("/api/clientID", func(c echo.Context) error {
		config, err := GetConfig()
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, config.ClientID)
	})
	e.GET("/api/branches", func(c echo.Context) error {
		branches, err := GetClient().Branches(context.TODO())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, branches)
	})
	e.GET("/api/connect", func(c echo.Context) error {
		config, err := GetConfig()
		if err != nil {
			return err
		}
		fmt.Println("connect")
		client := GetClient()
		err = client.CreateBranch(context.TODO(), config.ClientID, "测试分支")
		if err != nil {
			fmt.Println(err)
			return err
		}
		hash, err := client.WriteObject(context.TODO(), []byte("111"))
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = client.ReadObject(context.TODO(), hash, func(buf []byte) error {
			fmt.Println("ReadObject", string(buf))
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return c.String(http.StatusOK, "ok")
	})
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	e.Logger.Fatal(e.StartTLS(":"+port, "localhost.pem", "localhost-key.pem"))
}
