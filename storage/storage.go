package storage

import (
	"io"

	"github.com/lazyxu/kfs/kfscrypto"
)

const (
	TypBlob = iota
	TypTree
)

type Storage interface {
	Read(typ int, key string) (io.Reader, error)
	Write(typ int, reader io.Reader) (string, error)
	Delete(typ int, key string) error
}

type BaseStorage struct {
	hashFunc     func() kfscrypto.Hash
	checkOnWrite bool
	checkOnRead  bool
}

func NewBase(hashFunc func() kfscrypto.Hash, checkOnWrite bool, checkOnRead bool) BaseStorage {
	return BaseStorage{
		hashFunc:     hashFunc,
		checkOnRead:  checkOnWrite,
		checkOnWrite: checkOnRead,
	}
}
func (s *BaseStorage) HashFunc() kfscrypto.Hash {
	return s.hashFunc()
}

func (s *BaseStorage) CheckOnRead() bool {
	return s.checkOnRead
}

func (s *BaseStorage) CheckOnWrite() bool {
	return s.checkOnWrite
}
