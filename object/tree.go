package object

import (
	"bytes"

	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/storage"
)

type Tree struct {
	base  *Obj
	Items []*Metadata
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
	err := o.base.serializable.Serialize(o, b)
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
	return o.base.serializable.Deserialize(o, reader)
}

func (base *Obj) ReadDir(s storage.Storage, key string) (*Tree, error) {
	tree := base.NewTree()
	err := tree.Read(s, key)
	return tree, err
}

func (o *Tree) IsDir() bool {
	return true
}

func (o *Tree) IsFile() bool {
	return false
}
