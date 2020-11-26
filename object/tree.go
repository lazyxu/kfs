package object

import (
	"bytes"
	"encoding/gob"

	"github.com/lazyxu/kfs/storage/kfshash"

	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/storage"
)

type Tree struct {
	Items []*Metadata
}

var EmptyDir = &Tree{
	Items: make([]*Metadata, 0),
}
var EmptyDirHash string

func Init(hashFunc func() kfshash.Hash) error {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(EmptyDir)
	if err != nil {
		return err
	}
	EmptyDirHash, err = hashFunc().Cal(b)
	if err != nil {
		return err
	}
	EmptyFileHash, err = hashFunc().Cal(bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	return nil
}

func (o *Tree) GetNode(name string) (*Metadata, error) {
	for _, it := range o.Items {
		if it.Name == name {
			return it, nil
		}
	}
	return nil, e.ENoSuchFileOrDir
}

func (o *Tree) Write(s storage.Storage) (string, error) {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(o)
	if err != nil {
		return "", e.EWriteObject
	}
	return s.Write(storage.TypTree, b)
}

func (o *Tree) Read(s storage.Storage, key string) error {
	reader, err := s.Read(storage.TypTree, key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(reader).Decode(o)
}

func ReadDir(s storage.Storage, key string) (*Tree, error) {
	tree := new(Tree)
	err := tree.Read(s, key)
	return tree, err
}

func (o *Tree) IsDir() bool {
	return true
}

func (o *Tree) IsFile() bool {
	return false
}
