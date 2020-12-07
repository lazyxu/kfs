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
var serializable kfscrypto.Serializable
var hashFunc func() kfscrypto.Hash
var obj *object.Obj

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	hashFunc = func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	var err error
	s, err = fs.New("temp", hashFunc, true, true)
	if err != nil {
		panic(err)
	}
	serializable = &kfscrypto.GobEncoder{}
	obj = object.Init(hashFunc, serializable)
	e := echo.New()
	e.GET("/api/download/:hash", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		hash := c.Param("hash")
		b := obj.NewBlob()
		err := b.Read(s, hash)
		if err != nil {
			return err
		}
		c.Response().WriteHeader(http.StatusOK)
		_, err = io.Copy(c.Response(), b.Reader)
		if err != nil {
			return err
		}
		return nil
	})
	e.POST("/api/upload", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		b := obj.NewBlob()
		b.Reader = c.Request().Body
		hash, err := b.Write(s)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, hash)
	})
	e.Logger.Fatal(e.Start(":9999"))
}
