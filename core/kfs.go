package core

import (
	"github.com/lazyxu/kfs/dao"
	"os"
)

type KFS struct {
	Db           dao.Database
	S            dao.Storage
	thumbnailDir string
	transCodeDir string
	newStorage   func(root string) (dao.Storage, error)
	isSqlite     bool
}

func New(newDatabase func() (dao.Database, error), newStorage func() (dao.Storage, error)) (*KFS, error) {
	s, err := newStorage()
	if err != nil {
		return nil, err
	}
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s}, nil
}

func mkdirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(path, 0o700)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func New2(newDatabase func() (dao.Database, error), newStorage func() (dao.Storage, error), thumbnailDir, transCodeDir string) (*KFS, error) {
	s, err := newStorage()
	if err != nil {
		return nil, err
	}
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(thumbnailDir)
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(transCodeDir)
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s, thumbnailDir: thumbnailDir, transCodeDir: transCodeDir}, nil
}

func (fs *KFS) ThumbnailDir() string {
	return fs.thumbnailDir
}

func (fs *KFS) TransCodeDir() string {
	return fs.transCodeDir
}

func (fs *KFS) Close() error {
	err1 := fs.S.Close()
	err2 := fs.Db.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (fs *KFS) Reset() error {
	err := fs.S.Remove()
	if err != nil {
		return err
	}
	err = fs.S.Create()
	if err != nil {
		return err
	}
	err = fs.Db.Remove()
	if err != nil {
		return err
	}
	err = fs.Db.Create()
	if err != nil {
		return err
	}
	return nil
}
