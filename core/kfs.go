package core

import (
	"os"
	"path"

	"github.com/lazyxu/kfs/db/mysql"

	"github.com/lazyxu/kfs/db/gosqlite"

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

func New(root string) (*KFS, bool, error) {
	return NewWithSqlite(root, storage.NewStorage1)
}

func NewWithSqlite(root string, newStorage func(root string) (storage.Storage, error)) (*KFS, bool, error) {
	exist := true
	s, err := newStorage(root)
	if err != nil {
		return nil, false, err
	}
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
	return &KFS{Db: db, S: s, root: root, newStorage: newStorage, isSqlite: true}, exist, nil
}

func NewWithMysql(root string, newStorage func(root string) (storage.Storage, error)) (*KFS, error) {
	s, err := newStorage(root)
	if err != nil {
		return nil, err
	}
	db, err := mysql.Open("root:12345678@/kfs?charset=utf8&parseTime=true&multiStatements=true")
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s, root: root, newStorage: newStorage}, nil
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
