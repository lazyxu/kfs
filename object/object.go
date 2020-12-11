package object

import (
	"bytes"
	"io"

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
	r, err := o.EmptyDir.Serialize()
	if err != nil {
		panic(err)
	}
	o.EmptyDirHash, err = hw.Cal(r)
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
	return t.Write()
}

func (base *Obj) ReadTree(key string) (*Tree, error) {
	var t *Tree
	err := t.Read(key)
	if err != nil {
		return nil, err
	}
	return t, nil
}
