package dao

import "io"

type Storage interface {
	Write(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error)
	ReadWithSize(hash string) (SizedReadCloser, error)
	GetFilePath(hash string) string

	Remove() error
	Create() error
	Close() error
}

type SizedReadCloser interface {
	io.ReadSeekCloser
	Size() int64
}

func StorageNewFunc(root string, newStorage func(root string) (Storage, error)) func() (Storage, error) {
	return func() (Storage, error) {
		return newStorage(root)
	}
}
