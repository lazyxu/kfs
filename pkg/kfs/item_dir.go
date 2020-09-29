package kfs

import (
	"os"

	"github.com/lazyxu/kfs/kfs/e"
	"github.com/lazyxu/kfs/object"
)

type ItemDir struct {
	DItem
	object   *object.Dir
	children []Item
}

func newDir(kfs *KFS, name string) *ItemDir {
	return &ItemDir{
		DItem: DItem{
			kfs: kfs,
		},
		object:   object.NewItemDir(name),
		children: make([]Item, 0),
	}
}

func (i *ItemDir) Name() string {
	return i.object.Name()
}

func (i *ItemDir) IsDir() bool {
	return i.object.IsDir()
}

func (i *ItemDir) IsFile() bool {
	return i.object.IsFile()
}

func (i *ItemDir) loadObject() error {
	if i.object != nil {
		return nil
	}
	obj, err := i.kfs.scheduler.GetDirObjectByHash(i.object.Hash())
	if err != nil {
		return err
	}
	i.object = obj
	return nil
}

func (i *ItemDir) GetNode(name string) (Item, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err := i.loadObject()
	if err != nil {
		return nil, err
	}
	index, item := i.object.Get(name)
	if item == nil {
		return nil, e.ErrNotExist
	}
	// TODO: load children
	return i.children[index], nil
}

func (i *ItemDir) Add(item Item) error {
	err := i.loadObject()
	if err != nil {
		return err
	}
	exist := false
	i.object = i.object.Clone().(*object.Dir)
	i.object.Update(func(items []object.Object) []object.Object {
		for _, it := range items {
			if it.Name() == item.Name() {
				exist = true
				return items
			}
		}
		items = append(items, item.Item())
		return items
	})
	if exist {
		return e.ErrExist
	}
	i.children = append(i.children, item)
	return nil
}

func (i *ItemDir) Create(name string, flags int) (*ItemFile, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	f := newFile(i.kfs, name)
	err := i.Add(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (i *ItemDir) Remove(name string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err := i.loadObject()
	if err != nil {
		return err
	}
	i.object = i.object.Clone().(*object.Dir)
	index, err := i.object.Remove(name)
	if err != nil {
		return err
	}
	i.children = append(i.children[0:index], i.children[index+1:]...)
	return nil
}

func (i *ItemDir) ReadDirAll() ([]object.Object, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	err := i.loadObject()
	if err != nil {
		return nil, err
	}
	return i.object.Items(), nil
}

func (i *ItemDir) Mode() os.FileMode {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.object.Mode()
}

func (i *ItemDir) Item() object.Object {
	return i.object
}
