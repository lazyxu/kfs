package object

import (
	"bytes"
	"crypto/sha256"
	"io"

	"github.com/lazyxu/kfs/storage"
)

type Blob struct {
	Reader io.Reader
}

var EmptyFile = &Blob{
	Reader: bytes.NewReader([]byte{}),
}
var EmptyFileHash = string(sha256.New().Sum([]byte{}))

func (o *Blob) Write(s storage.Storage) (string, error) {
	return s.Write(storage.TypFile, o.Reader)
}

func (o *Blob) Read(s storage.Storage, key string) error {
	var err error
	o.Reader, err = s.Read(storage.TypFile, key)
	return err
}

func (o *Blob) IsDir() bool {
	return false
}

func (o *Blob) IsFile() bool {
	return true
}
