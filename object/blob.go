package object

import (
	"io"

	"github.com/lazyxu/kfs/storage"
)

type Blob struct {
	base   *BaseObject
	Reader io.Reader
}

func (o *Blob) Write(s storage.Storage) (string, error) {
	return s.Write(storage.TypBlob, o.Reader)
}

func (o *Blob) Read(s storage.Storage, key string) error {
	var err error
	o.Reader, err = s.Read(storage.TypBlob, key)
	return err
}

func (o *Blob) IsDir() bool {
	return false
}

func (o *Blob) IsFile() bool {
	return true
}
