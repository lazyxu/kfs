package object

import (
	"bytes"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage"
)

type Object interface {
	IsDir() bool
	IsFile() bool
	Write(s storage.Storage) (string, error)
	Read(s storage.Storage, key string) error
}

type BaseObject struct {
	serializable  kfscrypto.Serializable
	EmptyDirHash  string
	EmptyFileHash string
	EmptyFile     *Blob
	EmptyDir      *Tree
}

func Init(hashFunc func() kfscrypto.Hash, serializable kfscrypto.Serializable) *BaseObject {
	o := &BaseObject{serializable: serializable}
	o.EmptyFile = &Blob{
		base:   o,
		Reader: bytes.NewReader([]byte{}),
	}
	o.EmptyDir = &Tree{
		base:  o,
		Items: make([]*Metadata, 0),
	}
	hw := hashFunc()
	err := serializable.Serialize(o.EmptyDir, hw)
	if err != nil {
		panic(err)
	}
	o.EmptyDirHash, err = hw.Cal(nil)
	if err != nil {
		panic(err)
	}
	o.EmptyFileHash, err = hashFunc().Cal(o.EmptyFile.Reader)
	if err != nil {
		panic(err)
	}
	return o
}

func (base *BaseObject) NewBlob() *Blob {
	return &Blob{base: base}
}

func (base *BaseObject) NewTree() *Tree {
	return &Tree{base: base}
}
