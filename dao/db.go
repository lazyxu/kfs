package dao

import (
	"context"
	"os"
)

type Database interface {
	Remove() error
	Create() error
	Close() error

	ResetBranch(ctx context.Context, branchName string) error
	NewBranch(ctx context.Context, branchName string) (exist bool, err error)
	DeleteBranch(ctx context.Context, branchName string) error
	BranchInfo(ctx context.Context, branchName string) (branch Branch, err error)
	BranchList(ctx context.Context) (branches []IBranch, err error)

	WriteCommit(ctx context.Context, commit *Commit) error

	WriteDir(ctx context.Context, dirItems []DirItem) (dir Dir, err error)
	GetFile(ctx context.Context, branchName string, splitPath []string) (dirItem DirItem, err error)
	RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit Commit, branch Branch, err error)

	WriteFile(ctx context.Context, file File) error
	UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item DirItem) (commit Commit, branch Branch, err error)
	UpsertDirItems(ctx context.Context, branchName string, splitPath []string, items []DirItem) (commit Commit, branch Branch, err error)
	GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error)

	List(ctx context.Context, branchName string, splitPath []string) (dirItems []DirItem, err error)
	ListByHash(ctx context.Context, hash string) (dirItems []DirItem, err error)

	Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []DirItem, err error)
	Open2(ctx context.Context, branchName string, splitPath []string) (dirItem DirItem, dirItems []DirItem, err error)

	FileCount(ctx context.Context) (int, error)
	DirCount(ctx context.Context) (int, error)
	DirItemCount(ctx context.Context) (int, error)
	BranchCount(ctx context.Context) (int, error)
	UpsertDriverFile(ctx context.Context, f DriverFile) error
}

func DatabaseNewFunc(dataSourceName string, newDB func(dataSourceName string) (Database, error)) func() (Database, error) {
	return func() (Database, error) {
		return newDB(dataSourceName)
	}
}
