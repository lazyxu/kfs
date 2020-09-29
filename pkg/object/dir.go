package object

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/lazyxu/kfs/kfs/e"
)

type Dir struct {
	baseObject
	items []Object
}

func NewDir(items []Object) *Dir {
	return &Dir{items: items}
}

func (o *Dir) Items() []Object {
	items := make([]Object, len(o.items))
	copy(items, o.items)
	return items
}

func (o *Dir) Get(name string) (int, Object) {
	for i, it := range o.items {
		if it.Name() == name {
			return i, it.Clone()
		}
	}
	return 0, nil
}

func (o *Dir) Remove(name string) (int, error) {
	for i, it := range o.items {
		if it.Name() == name {
			o.items = append(o.items[:i], o.items[i+1:]...)
			o.hash = o.Hash()
			return i, nil
		}
	}
	return -1, e.ErrNotExist
}

func (o *Dir) Update(f func(items []Object) []Object) {
	o.items = f(o.Items())
	o.hash = o.Hash()
}

func (o *Dir) Hash() string {
	hash := sha256.New()
	hash.Write([]byte("dir"))
	for _, item := range o.items {
		binary.Write(hash, binary.LittleEndian, item.Name())
		binary.Write(hash, binary.LittleEndian, item.Hash())
		binary.Write(hash, binary.LittleEndian, item.BirthTime().UnixNano())
		binary.Write(hash, binary.LittleEndian, item.ModTime().UnixNano())
		binary.Write(hash, binary.LittleEndian, item.ChangeTime().UnixNano())
	}
	return string(hash.Sum(nil))
}

func (o Dir) Clone() Object {
	return &o
}
