package obj

import (
	"bytes"
	"crypto/sha256"
	"io"

	"github.com/lazyxu/kfs/storage"
	"github.com/lazyxu/kfs/storage/scheduler"
)

type File struct {
	Reader io.Reader
}

var EmptyFile = &File{
	Reader: bytes.NewReader([]byte{}),
}
var EmptyFileHash = string(sha256.New().Sum([]byte{}))

func (o *File) Write(s *scheduler.Scheduler) (string, error) {
	return s.WriteStream(storage.TypFile, o.Reader)
}

func (o *File) Read(s *scheduler.Scheduler, key string) error {
	var err error
	o.Reader, err = s.ReadStream(storage.TypFile, key)
	return err
}

func (o *File) IsDir() bool {
	return false
}

func (o *File) IsFile() bool {
	return true
}
