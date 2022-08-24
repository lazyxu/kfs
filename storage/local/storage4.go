package local

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gofrs/flock"
)

type Storage4 struct {
	root string
}

func NewStorage4(root string) (Storage, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, files), dirPerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	err = createGlobalLockFile(root)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &Storage4{root: root}, nil
}

var openedFiles [256]*os.File

func (s *Storage4) getFile(hash string) (*os.File, error) {
	id, err := strconv.ParseUint(hash[:2], 16, 8)
	if err != nil {
		return nil, err
	}
	f := openedFiles[id]
	if f != nil {
		return f, nil
	}
	filePath := path.Join(s.root, files, hash[:2])
	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		return nil, err
	}
	openedFiles[id] = f
	return f, nil
}

func (s *Storage4) WriteFn(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error) {
	lock := flock.New(path.Join(s.root, lockFileName))
	err := lock.Lock()
	if err != nil {
		return false, err
	}
	defer lock.Unlock()
	f, err := s.getFile(hash)
	lastOffset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	if err != nil {
		if os.IsPermission(err) {
			// file exists
			return true, nil
		}
		return false, err
	}
	hasher := sha256.New()
	err = fn(f, hasher)
	if err != nil {
		_, err = f.Seek(lastOffset, io.SeekStart)
		if err != nil {
			panic(err)
		}
		return false, err
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if hash != actual {
		_, err = f.Seek(lastOffset, io.SeekStart)
		if err != nil {
			panic(err)
		}
		return false, fmt.Errorf("invalid hash: expected %s, actual %s", hash, actual)
	}
	return false, nil
}

func (s *Storage4) ReadWithSize(hash string) (SizedReadCloser, error) {
	p := path.Join(s.root, files, hash[:2], hash[2:])
	f, err := os.OpenFile(p, os.O_RDONLY, 0o200)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return sizedReaderCloser{f, info.Size()}, nil
}
