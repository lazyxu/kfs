package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lazyxu/kfs/cmd/client/kfsclient"
	"github.com/lazyxu/kfs/cmd/client/pb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/api/connect", func(c echo.Context) error {
		//config, err := GetConfig()
		//if err != nil {
		//	return err
		//}
		fmt.Println("connect")
		client := kfsclient.New()
		status, err := client.Status(context.TODO(), &pb.Void{})
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(status)
		return c.String(http.StatusOK, "ok")
	})
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	e.Logger.Fatal(e.StartTLS(":"+port, "localhost.pem", "localhost-key.pem"))
}
