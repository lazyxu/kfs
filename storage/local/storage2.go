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

type Storage2 struct {
	root string
}

func NewStorage2(root string) (Storage, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Storage2{root: root}, nil
}

func (s *Storage2) getLocalLock(hash string) (*flock.Flock, error) {
	globalLock := flock.New(path.Join(s.root, lockFileName))
	err := globalLock.Lock()
	if err != nil {
		return nil, err
	}
	defer globalLock.Unlock()
	localLockFilePath := path.Join(s.root, files, hash[:2]+".lock")
	_, err = os.Stat(localLockFilePath)
	if os.IsNotExist(err) {
		lockFile, err := os.Create(localLockFilePath)
		if err != nil {
			return nil, err
		}
		defer lockFile.Close()
	} else if err != nil {
		return nil, err
	}
	return flock.New(localLockFilePath), nil
}

func (s *Storage2) WriteFn(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error) {
	localLock, err := s.getLocalLock(hash)
	if err != nil {
		return false, err
	}
	err = localLock.Lock()
	if err != nil {
		return false, err
	}
	defer localLock.Unlock()
	p := path.Join(s.root, files, hash)
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

func (s *Storage2) ReadWithSize(hash string) (SizedReadCloser, error) {
	p := path.Join(s.root, files, hash)
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

func (s *Storage2) Remove() error {
	return os.RemoveAll(s.root)
}

func (s *Storage2) Create() error {
	_, err := os.Stat(s.root)
	if err == nil {
		return fmt.Errorf("file or dir already exist: %s", s.root)
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(path.Join(s.root, files), dirPerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = createGlobalLockFile(s.root)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
