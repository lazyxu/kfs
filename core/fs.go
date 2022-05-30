package core

import (
	"context"
	"os"
	"path"

	storage "github.com/lazyxu/kfs/storage/local"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type FS interface {
	Checkout(ctx context.Context, branchName string) (bool, error)
	BranchInfo(ctx context.Context, branchName string) (branch sqlite.IBranch, err error)
	List(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item sqlite.IDirItem) error) error
	Upload(ctx context.Context, branchName string, dstPath string, srcPath string, uploadProcess UploadProcess) (sqlite.Commit, sqlite.Branch, error)
}

type KFS struct {
	Db *sqlite.DB
	S  *storage.Storage
}

func New(root string) (*KFS, bool, error) {
	s, err := storage.New(root)
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
		err = db.Reset()
		if err != nil {
			return nil, exist, err
		}
	}
	return &KFS{Db: db, S: s}, exist, nil
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
