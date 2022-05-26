package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

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

func (fs *KFS) Backup(ctx context.Context, root string, branchName string) error {
	backupCtx := storage.NewWalkerCtx[sqlite.FileOrDir](ctx, root, &uploadVisitor{fs: fs})
	ret, err := backupCtx.Scan()
	if err != nil {
		return err
	}
	if dir, ok := ret.(sqlite.Dir); ok {
		status := backupCtx.GetStatus()
		commit := sqlite.NewCommit(dir, branchName, fmt.Sprintf("%+v\n", status))
		err = fs.Db.WriteCommit(ctx, &commit)
		if err != nil {
			return err
		}
		branch := sqlite.NewBranch(branchName, commit, dir)
		err = fs.Db.WriteBranch(ctx, branch)
		if err != nil {
			return err
		}
	} else {
		return errors.New("expected a directory")
	}
	return nil
}

func FormatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	if splitPath[0] == "" {
		splitPath = splitPath[1:]
	}
	return splitPath
}

func (fs *KFS) BranchNew(ctx context.Context, branchName string) (bool, error) {
	return fs.Db.NewBranch(ctx, branchName)
}

func (fs *KFS) BranchInfo(ctx context.Context, branchName string) (branch sqlite.Branch, err error) {
	return fs.Db.BranchInfo(ctx, branchName)
}

func (fs *KFS) List(ctx context.Context, branchName string, p string) ([]sqlite.DirItem, error) {
	return fs.Db.List(ctx, branchName, FormatPath(p))
}

func (fs *KFS) Remove(ctx context.Context, branchName string, splitPath ...string) (sqlite.Commit, sqlite.Branch, error) {
	return fs.Db.Remove(ctx, branchName, splitPath)
}

func (fs *KFS) Cat(ctx context.Context, branchName string, p string) (io.ReadCloser, error) {
	hash, err := fs.Db.GetFileHash(ctx, branchName, FormatPath(p))
	if err != nil {
		return nil, err
	}
	return fs.S.Read(hash)
}

func (fs *KFS) Close() error {
	return fs.Db.Close()
}
