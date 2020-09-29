package kfs

import (
	"sync"

	"github.com/lazyxu/kfs/object"
)

type Item interface {
	Name() string
	IsDir() bool
	IsFile() bool
	Item() object.Object
	SetParent(dir *ItemDir)
}

type DItem struct {
	kfs    *KFS
	parent *ItemDir
	mutex  sync.RWMutex
}

func (i *DItem) SetParent(dir *ItemDir) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	i.parent = dir
}
