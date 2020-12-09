package fs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync/atomic"

	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/storage"
)

type Storage struct {
	storage.BaseStorage
	root string

	tempFileID uint32
}

const (
	dirPerm  = 0755
	filePerm = 0644
)

func typeToString(typ int) string {
	switch typ {
	case storage.TypBlob:
		return "blob"
	case storage.TypTree:
		return "tree"
	}
	return "unknown"
}

func (s *Storage) objectPath(typ int, key string) string {
	return path.Join(s.root, "objects", typeToString(typ), key)
}

func New(root string, hashFunc func() kfscrypto.Hash) (*Storage, error) {
	err := os.MkdirAll(path.Join(root, "objects", "tree"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "objects", "blob"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "refs"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "temp"), dirPerm)
	if err != nil {
		return nil, err
	}
	return &Storage{
		BaseStorage: storage.NewBase(hashFunc),
		root:        "temp",
		tempFileID:  0,
	}, nil
}

func (s *Storage) Read(typ int, key string) (io.Reader, error) {
	return os.Open(s.objectPath(typ, key))
}

func (s *Storage) Write(typ int, reader io.Reader) (string, error) {
	id := atomic.AddUint32(&s.tempFileID, 1)
	pTemp := path.Join(s.root, "temp", strconv.FormatUint(uint64(id), 10))
	fTemp, err := os.OpenFile(pTemp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return "", err
	}
	w := bufio.NewWriter(fTemp)
	hw := s.HashFunc()
	rr := io.TeeReader(reader, hw)
	_, err = w.ReadFrom(rr)
	if err != nil {
		fTemp.Close()
		return "", err
	}
	fTemp.Close()
	key, err := hw.Cal(nil)
	if err != nil {
		return "", err
	}
	p := s.objectPath(typ, key)
	fCurrent, err := os.OpenFile(p, os.O_RDONLY, filePerm)
	if err == nil {
		return key, fCurrent.Close()
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	// file not exists
	err = os.Rename(pTemp, p)
	if err != nil {
		return "", err
	}
	err = os.Chmod(p, 0444) // read only
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *Storage) Delete(typ int, key string) error {
	p := s.objectPath(typ, key)
	return os.Remove(p)
}

func (s *Storage) UpdateRef(name string, expect string, desire string) error {
	// TODO: expect
	return ioutil.WriteFile(path.Join(s.root, "refs", name), []byte(desire), filePerm)
}

func (s *Storage) GetRef(name string) (string, error) {
	bytes, err := ioutil.ReadFile(path.Join(s.root, "refs", name))
	return string(bytes), err
}
