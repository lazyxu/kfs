package node

import (
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/object"
)

type Node interface {
	os.FileInfo
	BTime() time.Time
	CTime() time.Time
	IsFile() bool
	Chmod(mode os.FileMode) error
	Stat() (os.FileInfo, error)
	ReadAt(buff []byte, off int64) (int, error)
	WriteAt(content []byte, offset int64) (n int, err error)
	Readdirnames(n int, offset int) (names []string, err error)
	Readdir(n int, offset int) ([]*object.Metadata, error)
	Close() error
	Path() string
	Truncate(size int64) error
	SetATime(t time.Time)
	SetMTime(t time.Time)
	Obj() *object.Obj
	Storage() storage.Storage
}

type ItemBase struct {
	obj      *object.Obj
	storage  storage.Storage
	Parent   *Dir
	mutex    sync.RWMutex
	metadata *object.Metadata
	dirty    uint64
}

func (i *ItemBase) Obj() *object.Obj {
	return i.obj
}

func (i *ItemBase) Storage() storage.Storage {
	return i.storage
}

func (i *ItemBase) Hash() string {
	return i.metadata.Hash()
}

func (i *ItemBase) Metadata() *object.Metadata {
	return i.metadata
}

func (i *ItemBase) BTime() time.Time {
	return i.metadata.BirthTime()
}

func (i *ItemBase) CTime() time.Time {
	return i.metadata.ChangeTime()
}

func (i *ItemBase) SetATime(t time.Time) {
	i.metadata = i.metadata.Builder().ChangeTime(t.UnixNano()).Build()
}

func (i *ItemBase) SetMTime(t time.Time) {
	i.metadata = i.metadata.Builder().ModifyTime(t.UnixNano()).Build()
}

func (i *ItemBase) Name() string {
	return i.metadata.Name()
}

func (i *ItemBase) Path() string {
	parent := i.Parent
	p := i.Name()
	for parent != nil {
		p = path.Join(parent.Name(), p)
		parent = parent.Parent
	}
	return "/" + p
}

func (i *ItemBase) Size() int64 {
	return i.metadata.Size()
}

func (i *ItemBase) Mode() os.FileMode {
	return i.metadata.Mode() & os.ModePerm
}

func (i *ItemBase) Chmod(mode os.FileMode) error {
	mode = (i.metadata.Mode() & (^os.ModePerm)) | (mode & os.ModePerm)
	i.metadata = i.metadata.Builder().Mode(mode).Build()
	return i.update()
}

func (i *ItemBase) ModTime() time.Time {
	return i.metadata.ModifyTime()
}

func (i *ItemBase) Sys() interface{} {
	return i.metadata
}

func (i *ItemBase) updateObj(o *object.Tree) error {
	hash, err := i.obj.WriteTree(o)
	if err != nil {
		return err
	}
	i.metadata = i.metadata.Builder().
		Hash(hash).Build()
	return i.update()
}

func (i *ItemBase) update() error {
	atomic.AddUint64(&i.dirty, 1)
	if i.Parent == nil {
		return nil
	}
	dd, err := i.Parent.load()
	if err != nil {
		return err
	}
	for index, item := range dd.Items {
		if item.Name() == i.metadata.Name() {
			items := append(dd.Items[0:index], i.metadata)
			dd.Items = append(items, dd.Items[index+1:]...)
			return i.Parent.updateObj(dd)
		}
	}
	return nil
}

func (i *ItemBase) IsDir() bool {
	return i.metadata.IsDir()
}

func (i *ItemBase) IsFile() bool {
	return i.metadata.IsFile()
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
