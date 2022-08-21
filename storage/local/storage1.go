package local

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/gofrs/flock"
)

type Storage1 struct {
	root string
}

func NewStorage1(root string) (*Storage1, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, files), dirPerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	err = createLockFile(root)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &Storage1{root: root}, nil
}

func (s *Storage1) WriteFn(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error) {
	lock := flock.New(path.Join(s.root, lockFileName))
	lock.Lock()
	defer lock.Unlock()
	dirPath := path.Join(s.root, files, hash[:2])
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, dirPerm)
		if err != nil {
			return false, err
		}
	} else if err != nil {
		return false, err
	}
	p := path.Join(dirPath, hash[2:])
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		if os.IsPermission(err) {
			// file exists
			return true, nil
		}
		return false, err
	}
	defer f.Close()
	hasher := sha256.New()
	err = fn(f, hasher)
	if err != nil {
		os.Remove(p)
		return false, err
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if hash != actual {
		os.Remove(p)
		return false, fmt.Errorf("invalid hash: expected %s, actual %s", hash, actual)
	}
	err = os.Chmod(p, 0o400) // read only
	if err != nil {
		os.Remove(p)
		return false, fmt.Errorf("failed to change file mode: %s", hash)
	}
	return false, nil
}

func (s *Storage1) ReadWithSize(hash string) (SizedReadCloser, error) {
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
