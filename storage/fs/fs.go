package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"sync/atomic"

	"github.com/lazyxu/kfs/storage/kfshash"

	"github.com/lazyxu/kfs/storage"
)

type Storage struct {
	root         string
	hashFunc     func() kfshash.Hash
	checkOnRead  bool
	checkOnWrite bool

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

func New(root string, hashFunc func() kfshash.Hash, checkOnRead bool, checkOnWrite bool) (*Storage, error) {
	err := os.MkdirAll(path.Join(root, "objects", "tree"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "objects", "blob"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "temp"), dirPerm)
	if err != nil {
		return nil, err
	}
	return &Storage{
		root:         root,
		hashFunc:     hashFunc,
		checkOnRead:  checkOnRead,
		checkOnWrite: checkOnWrite,
	}, nil
}

func (s *Storage) Read(typ int, key string) (io.Reader, error) {
	return os.Open(s.objectPath(typ, key))
}

func calc(f *os.File, h kfshash.Hash) (string, error) {
	return h.Cal(f)
}

func (s *Storage) Write(typ int, reader io.Reader) (string, error) {
	id := atomic.AddUint32(&s.tempFileID, 1)
	pTemp := path.Join(s.root, "temp", strconv.FormatUint(uint64(id), 10))
	fTemp, err := os.OpenFile(pTemp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return "", err
	}
	w := bufio.NewWriter(fTemp)
	_, err = w.ReadFrom(reader)
	if err != nil {
		fTemp.Close()
		return "", err
	}
	fTemp.Close()
	fTemp2, err := os.OpenFile(pTemp, os.O_RDONLY, filePerm)
	if err != nil {
		return "", err
	}
	defer fTemp2.Close()
	key, err := calc(fTemp2, s.hashFunc())
	if err != nil {
		return "", err
	}
	p := s.objectPath(typ, key)
	fCurrent, err := os.OpenFile(p, os.O_RDONLY, filePerm)
	if os.IsNotExist(err) {
		goto moveTempFile
	}
	if err != nil {
		return "", err
	}
	if s.checkOnWrite {
		actualKey, err := calc(fCurrent, s.hashFunc())
		if err != nil {
			return "", err
		}
		if actualKey != key {
			fmt.Fprintf(os.Stderr, "invalid object: expected %s, actual %s", key, actualKey)
			goto moveTempFile
		}
	}
	return key, nil
moveTempFile:
	err = os.Rename(pTemp, p)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *Storage) Delete(typ int, key string) error {
	p := path.Join(s.root, "objects", typeToString(typ), key)
	return os.Remove(p)
}
