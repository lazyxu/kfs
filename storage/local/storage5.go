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

	"github.com/lazyxu/kfs/dao"

	"github.com/gofrs/flock"
)

type Storage5 struct {
	root        string
	openedFiles [256]*os.File
}

func NewStorage5(root string) (dao.Storage, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Storage5{root: root}, nil
}

func (s *Storage5) getFile(hash string) (*os.File, error) {
	id, err := strconv.ParseUint(hash[:2], 16, 8)
	if err != nil {
		return nil, err
	}
	f := s.openedFiles[id]
	if f != nil {
		return f, nil
	}
	filePath := path.Join(s.root, files, hash[:2])
	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		return nil, err
	}
	s.openedFiles[id] = f
	return f, nil
}

func (s *Storage5) getLocalLock(hash string) (*flock.Flock, error) {
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

func (s *Storage5) Write(hash string, fn func(w io.Writer, hasher io.Writer) error) (bool, error) {
	//lock, err := s.getLocalLock(hash)
	//if err != nil {
	//	return false, err
	//}
	//err = lock.Lock()
	//if err != nil {
	//	return false, err
	//}
	//defer lock.Unlock()
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

func (s *Storage5) ReadWithSize(hash string) (dao.SizedReadCloser, error) {
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

func (s *Storage5) Remove() error {
	return os.RemoveAll(s.root)
}

func (s *Storage5) Create() error {
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

func (s *Storage5) Close() error {
	for _, f := range s.openedFiles {
		if f != nil {
			_ = f.Close()
		}
	}
	return nil
}
