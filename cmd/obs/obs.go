package main

import (
	"crypto/sha256"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage"
	"github.com/lazyxu/kfs/storage/fs"
)

func NewOBS() storage.Storage {
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	s, err := fs.New("kfs_root", hashFunc, true, true)
	if err != nil {
		panic(err)
	}
	return s
}
