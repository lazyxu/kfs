package object

import (
	"os"
	"time"
)

type Metadata struct {
	mode       os.FileMode
	birthTime  int64
	modifyTime int64
	changeTime int64
	size       int64
	name       string
	hash       string
}

func (i *Metadata) Mode() os.FileMode {
	return i.mode
}

func (i *Metadata) BirthTime() time.Time {
	return time.Unix(0, i.birthTime)
}

func (i *Metadata) ModifyTime() time.Time {
	return time.Unix(0, i.modifyTime)
}

func (i *Metadata) ChangeTime() time.Time {
	return time.Unix(0, i.changeTime)
}

func (i *Metadata) Size() int64 {
	return i.size
}

func (i *Metadata) Name() string {
	return i.name
}

func (i *Metadata) Hash() string {
	return i.hash
}

func (i *Metadata) IsFile() bool {
	return i.mode&S_IFREG != 0
}

func (i *Metadata) IsDir() bool {
	return i.mode&S_IFDIR != 0
}

func (base *Obj) NewDirMetadata(name string, perm os.FileMode) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		mode:       S_IFDIR | (perm & os.ModePerm),
		birthTime:  now,
		modifyTime: now,
		changeTime: now,
		name:       name,
		size:       0,
		hash:       base.EmptyDirHash,
	}
}

func (base *Obj) NewFileMetadata(name string, perm os.FileMode) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		mode:       S_IFREG | (perm & os.ModePerm),
		birthTime:  now,
		modifyTime: now,
		changeTime: now,
		name:       name,
		size:       0,
		hash:       base.EmptyFileHash,
	}
}
