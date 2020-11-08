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
	Chmod(mode os.FileMode) error
	Stat() (os.FileInfo, error)
	Read(buff []byte) (int, error)
	ReadAt(buff []byte, off int64) (int, error)
	Write(content []byte) (n int, err error)
	WriteAt(content []byte, offset int64) (n int, err error)
	Readdirnames(n int) (names []string, err error)
	Readdir(n int) ([]*object.Metadata, error)
	Close() error
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

func (i *ItemBase) Chmod(mode os.FileMode) error {
	metadata := *i.Metadata
	metadata.Mode = mode
	return i.update()
}

func (i *ItemBase) ModTime() time.Time {
	return time.Unix(0, i.ModifyTime)
}

func (i *ItemBase) Sys() interface{} {
	return i.Metadata
}

func (i *ItemBase) updateObj(o object.Object) error {
	hash, err := o.Write(i.kfs.scheduler)
	if err != nil {
		return err
	}
	i.Metadata.Hash = hash
	return i.update()
}

func (i *ItemBase) update() error {
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
			return i.parent.updateObj(dd)
		}
	}
	return nil
}

func (i *ItemBase) Stat() (os.FileInfo, error) {
	return i, nil
}

func (i *ItemBase) Close() error {
	if i == nil {
		return os.ErrInvalid
	}
	return i.update()
}
