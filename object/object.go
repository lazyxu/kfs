package object

import (
	"bytes"
	"io"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage"
)

type Object interface {
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

func (base *Obj) WriteBlob(r io.Reader) (string, error) {
	return base.s.Write(storage.TypBlob, r)
}

func (base *Obj) ReadBlob(key string) (io.Reader, error) {
	return base.s.Read(storage.TypBlob, key)
}

func (base *Obj) WriteTree(t *Tree) (string, error) {
	b := &bytes.Buffer{}
	err := serializable.Serialize(t, b)
	if err != nil {
		return "", e.EWriteObject
	}
	return base.s.Write(storage.TypTree, b)
}

func (base *Obj) ReadTree(key string) (*Tree, error) {
	reader, err := base.s.Read(storage.TypTree, key)
	if err != nil {
		return nil, err
	}
	var t *Tree
	err = serializable.Deserialize(t, reader)
	if err != nil {
		return nil, err
	}
	return t, nil
}
