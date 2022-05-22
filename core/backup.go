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
		commit := sqlite.NewCommit(dir, branchName, fmt.Sprintf("%+v\n", status))
		err = fs.db.WriteCommit(ctx, &commit)
		if err != nil {
			return err
		}
		branch := sqlite.NewBranch(branchName, commit, dir)
		err = fs.db.WriteBranch(ctx, branch)
		if err != nil {
			return err
		}
	} else {
		return errors.New("expected a directory ")
	}
	return nil
}

func formatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	if splitPath[0] == "" {
		splitPath = splitPath[1:]
	}
	return splitPath
}

func (fs *KFS) BranchNew(ctx context.Context, branchName string) (bool, error) {
	return fs.db.NewBranch(ctx, branchName)
}

func (fs *KFS) BranchInfo(ctx context.Context, branchName string) (branch sqlite.Branch, err error) {
	return fs.db.BranchInfo(ctx, branchName)
}

func (fs *KFS) Upload(ctx context.Context, fn func(f io.Writer, hasher io.Writer) error, branchName string, p string,
	hash string, size uint64, mode uint64, createTime uint64,
	modifyTime uint64, changeTime uint64, accessTime uint64) (exist bool, commit sqlite.Commit, err error) {
	exist, err = fs.s.WriteFn(hash, fn)
	if err != nil {
		return
	}
	commit, err = fs.db.UploadFile(ctx, branchName, formatPath(p), hash,
		size, mode, createTime, modifyTime, changeTime, accessTime)
	return
}

func (fs *KFS) List(ctx context.Context, branchName string, p string) ([]sqlite.DirItem, error) {
	return fs.db.List(ctx, branchName, formatPath(p))
}

func (fs *KFS) Remove(ctx context.Context, branchName string, splitPath ...string) (sqlite.Commit, error) {
	return fs.db.Remove(ctx, branchName, splitPath)
}

func (fs *KFS) Cat(ctx context.Context, branchName string, p string) (io.ReadCloser, error) {
	hash, err := fs.db.GetFileHash(ctx, branchName, formatPath(p))
	if err != nil {
		return nil, err
	}
	return fs.s.Read(hash)
}

func (fs *KFS) Close() error {
	return fs.db.Close()
}
