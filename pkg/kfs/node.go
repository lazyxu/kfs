package kfs

import (
	"os"
	"sync"

	"github.com/lazyxu/kfs/storage/obj"
)

type Node interface {
	Name() string
	IsDir() bool
	IsFile() bool
	GetMetadata() obj.Metadata
}

type ItemBase struct {
	kfs    *KFS
	parent *Dir
	mutex  sync.RWMutex
	obj.Metadata
}

func (i *ItemBase) GetMetadata() obj.Metadata {
	return i.Metadata
}

func (i *ItemBase) Mode() os.FileMode {
	return i.Metadata.Mode
}

func (i *ItemBase) Name() string {
	return i.Metadata.Name
}

func (i *ItemBase) Size() int64 {
	return i.Metadata.Size
}

func (i *ItemBase) update(o obj.Object) error {
	hash, err := o.Write(i.kfs.scheduler)
	if err != nil {
		return err
	}
	i.Metadata.Hash = hash
	if i.parent == nil {
		return nil
	}
	dd, err := i.parent.load()
	if err != nil {
		return err
	}
	for index, item := range dd.Items {
		if item.Name == i.Metadata.Name {
			items := append(dd.Items[0:index], i.Metadata)
			dd.Items = append(items, dd.Items[index+1:]...)
			return i.parent.update(dd)
		}
	}
	return nil
}
