package storage

import (
	"io"

	"github.com/lazyxu/kfs/storage/kfshash"
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
	hashFunc     func() kfshash.Hash
	checkOnWrite bool
	checkOnRead  bool
}

func NewBase(hashFunc func() kfshash.Hash, checkOnWrite bool, checkOnRead bool) BaseStorage {
	return BaseStorage{
		hashFunc:     hashFunc,
		checkOnRead:  checkOnWrite,
		checkOnWrite: checkOnRead,
	}
}
func (s *BaseStorage) HashFunc() kfshash.Hash {
	return s.hashFunc()
}

func (s *BaseStorage) CheckOnRead() bool {
	return s.checkOnRead
}

func (s *BaseStorage) CheckOnWrite() bool {
	return s.checkOnWrite
}
