package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"io/fs"
	"path/filepath"
)

func diskUsage() error {
	info, err := disk.Usage("C://")
	if err != nil {
		return err
	}
	fmt.Println("剩余空间", humanize.IBytes(info.Free))
	var thumbnailSize uint64
	err = filepath.Walk("thumbnail", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			thumbnailSize += uint64(info.Size())
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Println("缩略图总大小", humanize.IBytes(thumbnailSize))
	dbSize, err := kfsCore.Db.Size()
	if err != nil {
		return err
	}
	fmt.Println("元数据总大小", humanize.IBytes(uint64(dbSize)))
	fileSize, err := kfsCore.Db.SumFileSize(context.TODO())
	if err != nil {
		return err
	}
	fmt.Println("文件总大小", humanize.IBytes(fileSize))
	return nil
}
