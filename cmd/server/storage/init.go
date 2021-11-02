package storage

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/lazyxu/kfs/cmd/server/kfsserver/errorutil"

	"github.com/lazyxu/kfs/cmd/server/kfscrypto"
)

type Storage struct {
	root     string
	HashFunc func() kfscrypto.Hash
}

var (
	EmptyDir      = Directory{Items: make([]*Metadata, 0)}
	EmptyDirHash  string
	EmptyFileHash string
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

func mkdir(path string) {
	err := os.MkdirAll(path, dirPerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func New(rootDir string, hashFunc func() kfscrypto.Hash) (*Storage, error) {
	s := &Storage{root: rootDir, HashFunc: hashFunc}
	root, err := filepath.Abs(rootDir)
	errorutil.PanicIfErr(err)
	mkdir(path.Join(root, "branch"))
	mkdir(path.Join(root, "object"))
	println(hex.EncodeToString(s.HashFunc().Cal(nil)))
	buffer := new(bytes.Buffer)
	directoryEncoderDecoder.Encode(&EmptyDir, buffer)
	hw := s.HashFunc()
	hash := hw.Cal(buffer)
	directoryEncoderDecoder.Encode(&EmptyDir, buffer)
	s.WriteObject(hash, func(f func(reader io.Reader)) {
		f(buffer)
	})
	EmptyDirHash = hex.EncodeToString(hash)
	return s, nil
}
