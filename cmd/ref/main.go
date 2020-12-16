package main

import (
	"crypto/sha256"
	"net/http"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage/fs"
	"github.com/sirupsen/logrus"
)

var s storage.Storage
var hashFunc func() kfscrypto.Hash

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
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/api/ref/:name", func(c echo.Context) error {
		hash, err := s.GetRef(c.Param("name"))
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, hash)
	})
	e.PUT("/api/ref/:name/:old/:new", func(c echo.Context) error {
		err := s.UpdateRef(c.Param("name"), c.Param("old"), c.Param("new"))
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(":9998"))
}
