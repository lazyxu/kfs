package main

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/lazyxu/kfs/kfscore/kfscrypto"
	"github.com/lazyxu/kfs/kfscore/storage"
	"github.com/lazyxu/kfs/kfscore/storage/fs"
	"github.com/lazyxu/kfs/kfscore/storage/memory"

	"github.com/spf13/viper"
)

func initStorage() (s storage.Storage) {
	typ := viper.GetString("storage")
	fmt.Println("storage", typ)
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	if typ == "fileSystem" {
		root := viper.GetString("kfs-root")
		if root == "" {
			root = os.TempDir()
		}
		fmt.Println("root", root)
		var err error
		s, err = fs.New(root, hashFunc)
		if err != nil {
			panic(err)
		}
	} else if typ == "memory" {
		s = memory.New(hashFunc)
	}
	return
}
