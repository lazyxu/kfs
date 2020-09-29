package kfs

import (
	"os"
	"sync"
	"time"

	"github.com/lazyxu/kfs/kfs/e"

	"github.com/lazyxu/kfs/node"
)

type Item struct {
	hash   string
	parent *Dir
	node   node.Node
}

func (item *Item) Node() (node.Node, error) {
	return item.node, nil
}

type Dir struct {
	node.TimeImpl
	kfs   *KFS
	name  string
	path  string
	files map[string]*Item
	mutex sync.RWMutex // protects the following
}

func NewDir(kfs *KFS, name string) *Dir {
	now := time.Now()
	return &Dir{
		TimeImpl: node.TimeImpl{
			BTime: now,
			ATime: now,
			Mtime: now,
			CTime: now,
		},
		kfs:   kfs,
		path:  name,
		name:  name,
		files: make(map[string]*Item),
	}
}

// Stat looks up a specific entry in the receiver.
//
// Stat should return a Node corresponding to the entry.  If the
// name does not exist in the directory, Stat should return ErrNotExist.
//
// Stat need not to handle the names "." and "..".
func (d *Dir) Stat(leaf string) (*Item, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for name, item := range d.files {
		if name == leaf {
			return item, nil
		}
	}
	return nil, e.ErrNotExist
}

func (d *Dir) Add(name string, node node.Node) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.files[name]
	if ok {
		return e.ErrExist
	}
	d.files[name] = &Item{
		hash:   "",
		node:   node,
		parent: d,
	}
	return nil
}

func (d *Dir) Create(name string, flags int) (*File, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.files[name]
	if ok {
		return nil, e.ErrExist
	}
	f, err := NewFile(d.kfs, name)
	if err != nil {
		return nil, err
	}
	d.files[name] = &Item{
		hash:   "",
		node:   f,
		parent: d,
	}
	return f, nil
}

func (d *Dir) Remove(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.files[name]
	if !ok {
		return e.ErrNotExist
	}
	delete(d.files, name)
	return nil
}

func (d *Dir) ReadDirAll() (map[string]*Item, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.files, nil
}

func (d *Dir) Name() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.name
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) IsFile() bool {
	return false
}

func (d *Dir) Size() (int64, error) {
	return 0, nil
}

func (d *Dir) Mode() (mode os.FileMode) {
	return d.kfs.Opt.DirPerms
}

func (d *Dir) Truncate(size uint64) error {
	return e.ENotImpl
}
