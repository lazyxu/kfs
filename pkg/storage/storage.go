package storage

import (
	"github.com/lazyxu/kfs/object"
)

type Storage interface {
	Get(hash string) (object.Object, error)
	Add(obj object.Object) error
	Remove(hash string) error
}
