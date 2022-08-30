package local

import "io"

type Storage interface {
	WriteFn(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error)
	ReadWithSize(hash string) (SizedReadCloser, error)

	Remove() error
	Create() error
}

type SizedReadCloser interface {
	io.ReadCloser
	Size() int64
}

type sizedReaderCloser struct {
	io.ReadCloser
	size int64
}

func (rc sizedReaderCloser) Size() int64 {
	return rc.size
}
