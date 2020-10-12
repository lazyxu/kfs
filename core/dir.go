package core

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
)

type Dir struct {
	ItemBase
	items map[string]Node
}

func NewDir(kfs *KFS, name string, perm os.FileMode) *Dir {
	return &Dir{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: object.NewDirMetadata(name, perm),
		},
		items: make(map[string]Node),
	}
}

func (i *Dir) load() (*object.Tree, error) {
	tree := new(object.Tree)
	err := tree.Read(i.kfs.scheduler, i.Metadata.Hash)
	return tree, err
}

func getSize(r io.Reader) (int64, error) {
	switch v := r.(type) {
	case *bytes.Buffer:
		return int64(v.Len()), nil
	case *bytes.Reader:
		return int64(v.Len()), nil
	case *strings.Reader:
		return int64(v.Len()), nil
	case *os.File:
		info, err := v.Stat()
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	default:
		return 0, errors.New("invalid type: " + reflect.TypeOf(r).String())
	}
}

func (i *Dir) add(metadata *object.Metadata, item object.Object) error {
	d, err := i.load()
	if err != nil {
		return err
	}
	for _, it := range d.Items {
		if it.Name == metadata.Name {
			return e.ErrExist
		}
	}

	if blob, ok := item.(*object.Blob); ok {
		size, err := getSize(blob.Reader)
		if err != nil {
			return err
		}
		metadata.Size = size
	}
	itemHash, err := item.Write(i.kfs.scheduler)
	if err != nil {
		return err
	}
	metadata.Hash = itemHash
	d.Items = append(d.Items, metadata)

	return i.updateObj(d)
}

func (i *Dir) Create(name string, flags int) (*File, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	f := &File{
		ItemBase: ItemBase{
			kfs:      i.kfs,
			parent:   i,
			Metadata: object.NewFileMetadata(name),
		},
	}
	err := i.add(f.Metadata, object.EmptyFile)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (i *Dir) get(name string) (*object.Metadata, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	for _, item := range d.Items {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, e.ErrNotExist
}

func (i *Dir) remove(name string, all bool) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	d, err := i.load()
	if err != nil {
		return err
	}
	for index, item := range d.Items {
		if item.Name == name {
			if all || item.IsFile() || item.Hash == object.EmptyDirHash {
				d.Items = append(d.Items[0:index], d.Items[index+1:]...)
				delete(i.items, name)
				return i.updateObj(d)
			}
			if item.IsDir() {
				return e.ENotEmpty
			}
		}
	}
	if all {
		return nil
	}
	return e.ErrNotExist
}

func (i *Dir) ReadDirAll() ([]*object.Metadata, error) {
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	return d.Items, nil
}
