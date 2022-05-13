package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type KFS struct {
	db *sqlite.DB
	s  *storage.Storage
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
	return &KFS{db: db, s: s}, exist, nil
}

func (fs *KFS) Backup(ctx context.Context, root string, branchName string) error {
	backupCtx := storage.NewBackupCtx[sqlite.FileOrDir](ctx, root, &uploadVisitor{fs: fs})
	ret, err := backupCtx.Scan()
	if err != nil {
		return err
	}
	if dir, ok := ret.(sqlite.Dir); ok {
		status := backupCtx.GetStatus()
		commit := sqlite.NewCommit(dir, branchName)
		err = fs.db.WriteCommit(ctx, &commit)
		if err != nil {
			return err
		}
		branch := sqlite.NewBranch(branchName, fmt.Sprintf("%+v\n", status), commit, dir)
		err = fs.db.WriteBranch(ctx, branch)
		if err != nil {
			return err
		}
	} else {
		return errors.New("expected a directory ")
	}
	return nil
}

func (fs *KFS) List(ctx context.Context, branchName string, splitPath ...string) ([]sqlite.DirItem, error) {
	return fs.db.List(ctx, branchName, splitPath)
}

func (fs *KFS) Remove(ctx context.Context, branchName string, splitPath ...string) (err error) {
	return fs.db.Remove(ctx, branchName, splitPath)
}

func (fs *KFS) Close() error {
	return fs.db.Close()
}