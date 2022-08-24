package core

import (
	"os"
	"path"

	"github.com/lazyxu/kfs/db/gosqlite"

	storage "github.com/lazyxu/kfs/storage/local"
)

type UploadConfig struct {
	Encoder       string
	UploadProcess UploadProcess
	Concurrent    int
	Verbose       bool
}

type KFS struct {
	Db         *gosqlite.DB
	S          storage.Storage
	root       string
	newStorage func(root string) (storage.Storage, error)
}

func New(root string) (*KFS, bool, error) {
	return NewWithStorage(root, storage.NewStorage1)
}

func NewWithStorage(root string, newStorage func(root string) (storage.Storage, error)) (*KFS, bool, error) {
	s, err := newStorage(root)
	if err != nil {
		return nil, false, err
	}
	exist := true
	dbFileName := path.Join(root, "kfs.db")
	_, err = os.Stat(dbFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, false, err
		}
		exist = false
	}
	db, err := gosqlite.Open(dbFileName)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		err = db.Create()
		if err != nil {
			return nil, exist, err
		}
	}
	return &KFS{Db: db, S: s, root: root, newStorage: newStorage}, exist, nil
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
