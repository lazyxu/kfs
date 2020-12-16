package main

import (
	"crypto/sha256"
	"os"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/object"
	"github.com/lazyxu/kfs/storage"
	"github.com/lazyxu/kfs/storage/fs"
)

var s storage.Storage
var obj *object.Obj

func Init() {
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	var err error
	s, err = fs.New("temp", hashFunc)
	if err != nil {
		panic(err)
	}
	obj = object.Init(s)
	_, err = s.GetRef("default")
	if err == nil {
		return
	}
	kfs := core.New(&kfscommon.Options{
		UID:       uint32(os.Getuid()),
		GID:       uint32(os.Getgid()),
		DirPerms:  object.S_IFDIR | 0755,
		FilePerms: object.S_IFREG | 0644,
	}, s)
	err = kfs.Storage().UpdateRef("default", "", kfs.Root().Hash())
	if err != nil {
		panic(err)
	}
}
