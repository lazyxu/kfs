package filesystem

import (
	"os"

	"github.com/lazyxu/kfs/pkg/objectstore/common"
)

type fileReader struct {
	common.Reader
	common.ReaderWriter
	file *os.File
}

func (r *fileReader) Write(data []byte) (int, error) {
	return r.file.Write(data)
}

func (r *fileReader) Read(packet *[]byte) (int, error) {
	return r.file.Read(*packet)
}

func (r *fileReader) Close() error {
	return r.file.Close()
}
