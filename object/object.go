package object

import (
	"bytes"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage"
)

type Object interface {
	IsDir() bool
	IsFile() bool
	Write() (string, error)
	Read(key string) error
}

type Obj struct {
	s             storage.Storage
	EmptyDirHash  string
	EmptyFileHash string
	EmptyFile     *Blob
	EmptyDir      *Tree
}

var serializable kfscrypto.Serializable

func init() {
	serializable = &kfscrypto.GobEncoder{}
}

func Init(s storage.Storage) *Obj {
	o := &Obj{s: s}
	o.EmptyFile = &Blob{
		base:   o,
		Reader: bytes.NewReader([]byte{}),
	}
	o.EmptyDir = &Tree{
		base:  o,
		Items: make([]*Metadata, 0),
	}
	hw := s.HashFunc()
	err := serializable.Serialize(o.EmptyDir, hw)
	if err != nil {
		panic(err)
	}
	o.EmptyDirHash, err = hw.Cal(nil)
	if err != nil {
		panic(err)
	}
	o.EmptyFileHash, err = s.HashFunc().Cal(o.EmptyFile.Reader)
	if err != nil {
		panic(err)
	}
	return o
}

func (base *Obj) NewBlob() *Blob {
	return &Blob{base: base}
}

func (base *Obj) NewTree() *Tree {
	return &Tree{base: base}
}
