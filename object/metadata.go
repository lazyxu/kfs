package object

import (
	"os"
	"time"
)

type Metadata struct {
	Mode       os.FileMode
	BirthTime  int64
	ModifyTime int64
	ChangeTime int64
	Name       string
	Size       int64
	Hash       string
}

func (i *Metadata) IsFile() bool {
	return i.Mode&S_IFREG != 0
}

func (i *Metadata) IsDir() bool {
	return i.Mode&S_IFDIR != 0
}

func NewDirMetadata(name string, perm os.FileMode) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		Mode:       0040000 | perm,
		BirthTime:  now,
		ModifyTime: now,
		ChangeTime: now,
		Name:       name,
		Size:       0,
		Hash:       EmptyDirHash,
	}
}

func NewFileMetadata(name string) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		Mode:       DefaultFileMode,
		BirthTime:  now,
		ModifyTime: now,
		ChangeTime: now,
		Name:       name,
		Size:       0,
		Hash:       EmptyFileHash,
	}
}
