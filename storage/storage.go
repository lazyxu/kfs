package storage

import (
	"io"

	"github.com/lazyxu/kfs/kfscrypto"
)

const (
	TypBlob = iota
	TypTree
	TypRef
)

type Storage interface {
	Read(typ int, key string) (io.Reader, error)
	Write(typ int, reader io.Reader) (string, error)
	Delete(typ int, key string) error
	UpdateRef(name string, expect string, desire string) error
	GetRef(name string) (string, error)
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
