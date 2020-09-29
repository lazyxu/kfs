package object

import (
	"os"
	"time"
)

type Object interface {
	BirthTime() time.Time
	ModTime() time.Time
	ChangeTime() time.Time

	SetHash(hash string)
	Hash() string
	Name() string
	Size() int64
	Mode() os.FileMode
	IsDir() bool
	IsFile() bool

	Clone() Object
}

type baseObject struct {
	TimeImpl
	name string
	hash string
	size int64
	mode os.FileMode
}

func (item *baseObject) SetHash(hash string) {
	item.hash = hash
	item.TimeImpl.MTime = time.Now()
}

func NewItemDir(name string) *Dir {
	return &Dir{
		baseObject: baseObject{
			TimeImpl: NewTimeImpl(),
			name:     name,
			hash:     EmptyDirHash,
			mode:     DefaultDirMode,
		},
		items: make([]Object, 0),
	}
}

func (item *baseObject) Hash() string {
	return item.hash
}

func (item *baseObject) Name() string {
	return item.name
}

func (item *baseObject) Size() int64 {
	return item.size
}

func (item *baseObject) Mode() os.FileMode {
	return item.mode
}

func (item *baseObject) IsDir() bool {
	return item.mode.IsDir()
}

func (item *baseObject) IsFile() bool {
	return item.mode.IsRegular()
}
