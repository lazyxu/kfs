package local

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
)

const (
	dirPerm      = 0o700
	lockFileName = "index.lock"
)

const (
	files = "files"
)

type Storage interface {
	WriteFn(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error)
	ReadWithSize(hash string) (SizedReadCloser, error)
}

func createGlobalLockFile(root string) error {
	lockFile, err := os.Create(path.Join(root, lockFileName))
	if err != nil {
		return err
	}
	defer lockFile.Close()
	return nil
}

func NewContent(str string) (string, []byte) {
	content := []byte(str)
	hasher := sha256.New()
	_, err := hasher.Write(content)
	if err != nil {
		panic(err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, content
}

func Write(s Storage, hash string, reader io.Reader) (bool, error) {
	return s.WriteFn(hash, func(f io.Writer, hasher io.Writer) error {
		rr := io.TeeReader(reader, hasher)
		_, err := io.Copy(f, rr)
		return err
	})
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
