package storage

import (
	"io"
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
