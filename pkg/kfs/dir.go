package kfs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/lazyxu/kfs/storage/obj"

	"github.com/lazyxu/kfs/kfs/e"
)

type Dir struct {
	ItemBase
	items map[string]Node
}

func NewDir(kfs *KFS, name string) *Dir {
	return &Dir{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: obj.NewDirMetadata(name),
		},
		items: make(map[string]Node),
	}
}

func (i *Dir) load() (*obj.Dir, error) {
	d := new(obj.Dir)
	err := d.Read(i.kfs.scheduler, i.Metadata.Hash)
	return d, err
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

func (i *Dir) Add(metadata obj.Metadata, item obj.Object) error {
	d, err := i.load()
	if err != nil {
		return err
	}
	for _, it := range d.Items {
		if it.Name == metadata.Name {
			return e.ErrExist
		}
	}

	if f, ok := item.(*obj.File); ok {
		size, err := getSize(f.Reader)
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

	return i.update(d)
}

func (i *Dir) Create(name string, flags int) (*File, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	f := &File{
		ItemBase: ItemBase{
			kfs:      i.kfs,
			parent:   i,
			Metadata: obj.NewFileMetadata(name),
		},
	}
	err := i.Add(f.Metadata, obj.EmptyFile)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (i *Dir) Remove(name string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	d, err := i.load()
	if err != nil {
		return err
	}
	for index, item := range d.Items {
		if item.Name == name {
			d.Items = append(d.Items[0:index], d.Items[index+1:]...)
			return i.update(d)
		}
	}
	return e.ErrNotExist
}

func (i *Dir) ReadDirAll() ([]obj.Metadata, error) {
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	return d.Items, nil
}
