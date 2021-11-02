package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"google.golang.org/grpc/status"

	"github.com/lazyxu/kfs/cmd/client/pb"

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

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusOK, &Response{
		Code:    code,
		Message: message,
	})
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, &Response{
		Code: 0,
		Data: data,
	})
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
	e.GET("/api/listBranches", func(c echo.Context) error {
		branches, err := GetClient().ListBranches(context.TODO())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, branches)
	})
	e.POST("/api/createBranch", func(c echo.Context) error {
		branch := &pb.Branch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.CreateBranch(context.TODO(), branch)
		if err != nil {
			errStatus, _ := status.FromError(err)
			return Error(c, int(errStatus.Code()), errStatus.Message())
		}
		return Success(c, nil)
	})
	e.POST("/api/deleteBranch", func(c echo.Context) error {
		branch := &pb.Branch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.DeleteBranch(context.TODO(), branch)
		if err != nil {
			errStatus, _ := status.FromError(err)
			return Error(c, int(errStatus.Code()), errStatus.Message())
		}
		return Success(c, nil)
	})
	e.POST("/api/renameBranch", func(c echo.Context) error {
		branch := &pb.RenameBranch{}
		err := c.Bind(branch)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = GetClient().PbClient.RenameBranch(context.TODO(), branch)
		if err != nil {
			errStatus, _ := status.FromError(err)
			return Error(c, int(errStatus.Code()), errStatus.Message())
		}
		return Success(c, nil)
	})
	e.GET("/api/connect", func(c echo.Context) error {
		fmt.Println("connect")
		client := GetClient()
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
