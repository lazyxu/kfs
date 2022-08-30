package dao

import "io"

type Storage interface {
	Write(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error)
	ReadWithSize(hash string) (SizedReadCloser, error)

	Remove() error
	Create() error
}

type SizedReadCloser interface {
	io.ReadCloser
	Size() int64
}
