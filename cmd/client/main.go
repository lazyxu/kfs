package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

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

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, this is kfs client!")
	})
	e.GET("/api/clientID", func(c echo.Context) error {
		file, err := ioutil.ReadFile("kfs-config.json")
		if err != nil {
			return err
		}
		data := Config{}
		err = json.Unmarshal(file, &data)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, data.ClientID)
	})
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	e.Logger.Fatal(e.StartTLS(":"+port, "localhost.pem", "localhost-key.pem"))
}
