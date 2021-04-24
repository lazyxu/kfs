package storage

import (
	"io"

	"github.com/lazyxu/kfs/kfscore/kfscrypto"
)

const (
	TypBlob = iota
	TypTree
	TypRef
)

type Status struct {
	TotalPhysicalSize uint64
	BlobLogicalSize   uint64
	BlobCount         uint64
	TreeCount         uint64
}

type Storage interface {
	Read(typ int, key string, f func(reader io.Reader) error) error
	Write(typ int, reader io.Reader) (string, error)
	Exist(typ int, key string) (bool, error)
	//Commit(typ int, key string) error
	Delete(typ int, key string) error
	UpdateRef(name string, expect string, desire string) error
	GetRef(name string) (string, error)
	GetRefs() ([]string, error)
	Status() (Status, error)
	HashFunc() kfscrypto.Hash
}

type BaseStorage struct {
	hashFunc func() kfscrypto.Hash
}

func NewBase(hashFunc func() kfscrypto.Hash) BaseStorage {
	return BaseStorage{
		hashFunc: hashFunc,
	}
}
func (s *BaseStorage) HashFunc() kfscrypto.Hash {
	return s.hashFunc()
}
