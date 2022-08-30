package core

import (
	"github.com/lazyxu/kfs/dao"

	storage "github.com/lazyxu/kfs/storage/local"
)

type UploadConfig struct {
	Encoder       string
	UploadProcess UploadProcess
	Concurrent    int
	Verbose       bool
}

type KFS struct {
	Db         dao.DB
	S          storage.Storage
	root       string
	newStorage func(root string) (storage.Storage, error)
	isSqlite   bool
}

func New(funcNewDb func() (dao.DB, error), funcNewStorage func() (storage.Storage, error)) (*KFS, error) {
	s, err := funcNewStorage()
	if err != nil {
		return nil, err
	}
	db, err := funcNewDb()
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s}, nil
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
