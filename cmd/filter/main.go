package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/lazyxu/kfs/pkg/objectstore/common"

	"github.com/lazyxu/kfs/pkg/objectstore/filesystem"

	"github.com/lazyxu/kfs/pkg/hashfunc"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/pkg/ignorewalker"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := readConfig()
	logrus.SetLevel(logrus.InfoLevel)

	dirIgnore, err := ignorewalker.Walk(conf.Walk)
	if err != nil {
		logrus.WithError(err).Error("walk")
	}
	fmt.Printf("files: %d\n", len(dirIgnore.Files))
	hashFunc := hashfunc.GetHashFunc(hashfunc.HASH_SHA256)
	store := filesystem.New(conf.Root, hashFunc, binary.LittleEndian)
	for _, file := range dirIgnore.Files {
		data, err := ioutil.ReadFile(file.Path)
		if err != nil {
			logrus.WithError(err).Error("ReadFile")
		}
		hash, err := hashFunc.Hash(data)
		if err != nil {
			logrus.WithError(err).Error("Hash")
		}
		success, err := store.WriteObject(common.TypeBlob, hash, data, []byte(file.Path))
		logrus.WithError(err).WithFields(logrus.Fields{
			"hash":    hex.EncodeToString(hash),
			"success": success,
		}).Info(file.Path)
	}
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
