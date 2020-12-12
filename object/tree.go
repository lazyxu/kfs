package object

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/storage"
)

type Tree struct {
	base  *Obj
	Items []*Metadata
}

func (o *Tree) GetNode(name string) (*Metadata, error) {
	for _, it := range o.Items {
		if it.name == name {
			return it, nil
		}
	}
	return nil, e.ENoSuchFileOrDir
}

func (o *Tree) Write() (string, error) {
	r, err := o.Serialize()
	if err != nil {
		return "", err
	}
	return o.base.s.Write(storage.TypTree, r)
}

func (o *Tree) Serialize() (io.Reader, error) {
	b := &bytes.Buffer{}
	var err error
	for _, item := range o.Items {
		err = binary.Write(b, binary.LittleEndian, uint32(item.mode))
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		err = binary.Write(b, binary.LittleEndian, item.birthTime)
		if err != nil {
			return nil, err
		}
		err = binary.Write(b, binary.LittleEndian, item.modifyTime)
		if err != nil {
			return nil, err
		}
		err = binary.Write(b, binary.LittleEndian, item.changeTime)
		if err != nil {
			return nil, err
		}
		err = binary.Write(b, binary.LittleEndian, item.size)
		if err != nil {
			return nil, err
		}
		_, err = b.Write([]byte(item.hash))
		if err != nil {
			return nil, err
		}
		_, err = b.WriteString(item.name)
		if err != nil {
			return nil, err
		}
		_, err = b.WriteString("\n")
		if err != nil {
			return nil, err
		}
	}
	return b, err
}

func (o *Tree) Deserialize(b io.Reader) error {
	var err error
	for {
		item := new(Metadata)
		err = binary.Read(b, binary.LittleEndian, &item.mode)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = binary.Read(b, binary.LittleEndian, &item.birthTime)
		if err != nil {
			return err
		}
		err = binary.Read(b, binary.LittleEndian, &item.modifyTime)
		if err != nil {
			return err
		}
		err = binary.Read(b, binary.LittleEndian, &item.changeTime)
		if err != nil {
			return err
		}
		err = binary.Read(b, binary.LittleEndian, &item.size)
		if err != nil {
			return err
		}
		hash := make([]byte, len(o.base.EmptyDirHash))
		_, err = b.Read(hash)
		if err != nil {
			return err
		}
		item.hash = string(hash)
		for {
			temp := make([]byte, 1)
			_, err = b.Read(temp)
			if err != nil {
				return err
			}
			if temp[0] == '\n' {
				break
			}
			item.name += string(temp)
		}
		o.Items = append(o.Items, item)
	}
}

func (o *Tree) Read(key string) error {
	b, err := o.base.s.Read(storage.TypTree, key)
	if err != nil {
		return err
	}
	return o.Deserialize(b)
}

func (base *Obj) ReadDir(s storage.Storage, key string) (*Tree, error) {
	tree := base.NewTree()
	err := tree.Read(key)
	return tree, err
}
