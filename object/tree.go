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

func (o *Tree) Write() (string, error) {
	b := &bytes.Buffer{}
	err := serializable.Serialize(o, b)
	if err != nil {
		return "", e.EWriteObject
	}
	return o.base.s.Write(storage.TypTree, b)
}

func (o *Tree) Read(key string) error {
	reader, err := o.base.s.Read(storage.TypTree, key)
	if err != nil {
		return err
	}
	return serializable.Deserialize(o, reader)
}

func (base *Obj) ReadDir(s storage.Storage, key string) (*Tree, error) {
	tree := base.NewTree()
	err := tree.Read(key)
	return tree, err
}

func (o *Tree) IsDir() bool {
	return true
}

func (o *Tree) IsFile() bool {
	return false
}
