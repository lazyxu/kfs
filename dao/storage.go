package dao

import "io"

type Storage interface {
	Write(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error)
	ReadWithSize(hash string) (SizedReadCloser, error)

	Remove() error
	Create() error
	Close() error
}

type SizedReadCloser interface {
	io.ReaderAt
	io.Seeker
	io.ReadCloser
	Size() int64
}

func StorageNewFunc(root string, newStorage func(root string) (Storage, error)) func() (Storage, error) {
	return func() (Storage, error) {
		return newStorage(root)
	}
}
