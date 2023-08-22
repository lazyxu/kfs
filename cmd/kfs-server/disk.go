package main

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/disk"
	"io/fs"
	"path/filepath"
)

type DiskUsage struct {
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Thumbnail uint64 `json:"thumbnail"`
	Metadata  uint64 `json:"metadata"`
	File      uint64 `json:"file"`
}

func apiDiskUsage(c echo.Context) error {
	if !kfsCore.Db.IsSqlite() {
		return errors.New("is not sqlite")
	}
	var usage DiskUsage
	abs, err := filepath.Abs(kfsCore.Db.DataSourceName())
	if err != nil {
		return err
	}
	info, err := disk.Usage(filepath.Dir(abs))
	if err != nil {
		return err
	}
	usage.Total = info.Total
	usage.Free = info.Free
	err = filepath.Walk("thumbnail", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			usage.Thumbnail += uint64(info.Size())
		}
		return nil
	})
	if err != nil {
		return err
	}
	dbSize, err := kfsCore.Db.Size()
	if err != nil {
		return err
	}
	usage.Metadata = uint64(dbSize)
	usage.File, err = kfsCore.Db.SumFileSize(context.TODO())
	if err != nil {
		return err
	}
	return ok(c, usage)
}
