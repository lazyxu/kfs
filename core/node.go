package core

import (
	"os"
	"sync"
	"time"

	"github.com/lazyxu/kfs/object"
)

type Node interface {
	os.FileInfo
	IsFile() bool
}

type ItemBase struct {
	kfs    *KFS
	parent *Dir
	mutex  sync.RWMutex
	*object.Metadata
}

func (i *ItemBase) Name() string {
	return i.Metadata.Name
}

func (i *ItemBase) Size() int64 {
	return i.Metadata.Size
}

func (i *ItemBase) Mode() os.FileMode {
	return i.Metadata.Mode
}

func (i *ItemBase) ModTime() time.Time {
	return time.Unix(0, i.ModifyTime)
}

func (i *ItemBase) Sys() interface{} {
	return i.Metadata
}

func (i *ItemBase) update(o object.Object) error {
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
