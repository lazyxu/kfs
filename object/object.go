package object

import (
	"github.com/lazyxu/kfs/storage"
)

type Object interface {
	IsDir() bool
	IsFile() bool
	Write(s storage.Storage) (string, error)
	Read(s storage.Storage, key string) error
}
