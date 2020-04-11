package main

import (
	"fmt"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/pkg/ignorewalker"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	dirIgnore, err := ignorewalker.Walk("../../../..")
	if err != nil {
		logrus.WithError(err).Error("walk")
	}
	fmt.Printf("files: %d\n", len(dirIgnore.Files))
	fmt.Printf("filesSize: %s\n", humanize.Bytes(dirIgnore.Size))
	for _, file := range dirIgnore.Files {
		if file.Size > 10*1000*1000 {
			fmt.Printf("%s: %s\n", humanize.Bytes(file.Size), file.Path)
		}
	}
	fmt.Printf("repos: %d\n", len(dirIgnore.Children))
	dirIgnore.CalcDirSize()
	fmt.Printf("dirs: %d\n", len(dirIgnore.DirSize))
	for p, size := range dirIgnore.DirSize {
		if size > 100*1000*1000 {
			fmt.Printf("%s: %s\n", humanize.Bytes(size), p)
		}
	}
	logrus.Info("done!!!")
}
