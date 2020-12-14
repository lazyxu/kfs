package main

import (
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage/fs"
	"github.com/sirupsen/logrus"
)

var s storage.Storage
var hashFunc func() kfscrypto.Hash
var obj *object.Obj

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	hashFunc = func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	var err error
	s, err = fs.New("temp", hashFunc)
	if err != nil {
		panic(err)
	}
	obj = object.Init(s)
	e := echo.New()
	e.GET("/api/download/:hash", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		hash := c.Param("hash")
		err := obj.ReadBlob(hash, func(r io.Reader) error {
			_, err := io.Copy(c.Response(), r)
			return err
		})
		if err != nil {
			return err
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})
	e.OPTIONS("/api/upload", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "*")
		c.Response().Header().Set("Access-Control-Expose-Headers", "Authorization")
		return nil
	})
	e.POST("/api/upload", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "*")
		c.Response().Header().Set("Access-Control-Expose-Headers", "Authorization")
		hash, err := obj.WriteBlob(c.Request().Body)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, hash)
	})
	e.Logger.Fatal(e.StartTLS(":9999", "cert.pem", "key.pem"))
}
