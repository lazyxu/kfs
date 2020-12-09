package object

import (
	"io"

	"github.com/lazyxu/kfs/storage"
)

type Blob struct {
	base   *Obj
	Reader io.Reader
}

func (o *Blob) Write() (string, error) {
	return o.base.s.Write(storage.TypBlob, o.Reader)
}

func (o *Blob) Read(key string) error {
	var err error
	o.Reader, err = o.base.s.Read(storage.TypBlob, key)
	return err
}
