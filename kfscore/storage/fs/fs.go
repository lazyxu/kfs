package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/lazyxu/kfs/kfscore/kfscrypto"

	"github.com/lazyxu/kfs/kfscore/storage"
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
	return path.Join(s.root, typeToString(typ), key)
}

func (s *Storage) stageObjectPath(typ int, key string) string {
	return path.Join(s.root, "stage", typeToString(typ), key)
}

func New(root string, hashFunc func() kfscrypto.Hash) (*Storage, error) {
	err := os.MkdirAll(path.Join(root, "tree"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "blob"), dirPerm)
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
	err = os.MkdirAll(path.Join(root, "stage", "tree"), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(root, "stage", "blob"), dirPerm)
	if err != nil {
		return nil, err
	}
	return &Storage{
		BaseStorage: storage.NewBase(hashFunc),
		root:        root,
		tempFileID:  0,
	}, nil
}

func (s *Storage) Read(typ int, key string, f func(reader io.Reader) error) error {
	file, err := os.Open(s.objectPath(typ, key))
	if err != nil {
		return err
	}
	defer file.Close()
	return f(file)
}

func (s *Storage) Write(typ int, reader io.Reader) (string, error) {
	id := atomic.AddUint32(&s.tempFileID, 1)
	pTemp := path.Join(s.root, "temp", strconv.FormatUint(uint64(id), 10))
	fTemp, err := os.OpenFile(pTemp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return "", err
	}
	hw := s.HashFunc()
	rr := io.TeeReader(reader, hw)
	_, err = io.Copy(fTemp, rr)
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

//func (s *Storage) Transaction(f func() error) (uint32, error) {
//	pCommit := path.Join(s.root, "stage", strconv.FormatUint(uint64(id), 10))
//	err := os.RemoveAll(pCommit)
//	if err != nil {
//		return id, err
//	}
//	err = os.MkdirAll(pCommit, dirPerm)
//	if err != nil {
//		return id, err
//	}
//	err = os.MkdirAll(path.Join(pCommit, "tree"), dirPerm)
//	if err != nil {
//		return id, err
//	}
//	err = os.MkdirAll(path.Join(pCommit, "blob"), dirPerm)
//	if err != nil {
//		return id, err
//	}
//	err = f()
//	if err != nil {
//		return id, err
//	}
//	return id, nil
//}

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
	if os.IsNotExist(err) {
		err = nil
	}
	return string(bytes), err
}

func (s *Storage) GetRefs() ([]string, error) {
	infos, err := ioutil.ReadDir(path.Join(s.root, "refs"))
	if err != nil {
		return nil, err
	}
	branches := make([]string, len(infos))
	for i, info := range infos {
		branches[i] = info.Name()
	}
	return branches, err
}

func (s *Storage) Status() (status storage.Status, err error) {
	err = filepath.Walk(s.root, func(path string, info os.FileInfo, err error) error {
		status.TotalPhysicalSize += uint64(info.Size())
		return nil
	})
	infos, err := ioutil.ReadDir(path.Join(s.root, "blob"))
	if err != nil {
		return status, err
	}
	for _, info := range infos {
		status.BlobLogicalSize += uint64(info.Size())
	}
	status.BlobCount = uint64(len(infos))
	infos, err = ioutil.ReadDir(path.Join(s.root, "tree"))
	if err != nil {
		return status, err
	}
	status.TreeCount = uint64(len(infos))
	return status, err
}
