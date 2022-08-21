package core

import (
	"os"
	"path"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type KFS struct {
	Db   *sqlite.DB
	S    storage.Storage
	root string
}

func New(root string) (*KFS, bool, error) {
	s, err := storage.NewStorage1(root)
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
	db, err := sqlite.Open(dbFileName)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		err = db.Create()
		if err != nil {
			return nil, exist, err
		}
	}
	return &KFS{Db: db, S: s, root: root}, exist, nil
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
