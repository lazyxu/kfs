package main

import (
	"fmt"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/spf13/viper"

	"github.com/lazyxu/kfs/warpper/cgofuse"
)

func initFuse(s storage.Storage) {
	lib := viper.GetString("fuse-lib")
	if lib == "cgofuse" {
		mountPoint := viper.GetString("fuse-mount-point")
		fmt.Println("mount by", lib, mountPoint)
		if s == nil {
			panic("storage is nil")
		}
		cgofuse.Start(mountPoint, s)
	}
}
