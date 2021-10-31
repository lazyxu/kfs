package storage

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"

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

func mkdir(path string) error {
	err := os.MkdirAll(path, dirPerm)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func New(rootDir string, hashFunc func() kfscrypto.Hash) (*Storage, error) {
	s := &Storage{root: rootDir, HashFunc: hashFunc}
	root, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	err = mkdir(path.Join(root, "branch"))
	if err != nil {
		return nil, err
	}
	err = mkdir(path.Join(root, "object"))
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	err = directoryEncoderDecoder.Encode(&EmptyDir, buffer)
	if err != nil {
		panic(err)
	}
	hw := s.HashFunc()
	hash, err := hw.Cal(buffer)
	if err != nil {
		panic(err)
	}
	err = s.WriteObject(hash, func(f func(reader io.Reader) error) error {
		return f(buffer)
	})
	EmptyDirHash = hex.EncodeToString(hash)
	if err != nil {
		panic(err)
	}
	return s, nil
}
