package dao

import (
	"context"
	"os"
)

type DB interface {
	Create() error
	Close() error

	WriteBranch(ctx context.Context, branch Branch) error
	NewBranch(ctx context.Context, branchName string) (exist bool, err error)
	BranchInfo(ctx context.Context, branchName string) (branch Branch, err error)

	WriteCommit(ctx context.Context, commit *Commit) error

	WriteDir(ctx context.Context, dirItems []DirItem) (dir Dir, err error)
	GetFileHash(ctx context.Context, branchName string, splitPath []string) (hash string, err error)
	Remove(ctx context.Context, branchName string, splitPath []string) (commit Commit, branch Branch, err error)

	WriteFile(ctx context.Context, file File) error
	UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item DirItem) (commit Commit, branch Branch, err error)
	GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error)

	List(ctx context.Context, branchName string, splitPath []string) (dirItems []DirItem, err error)
	ListByHash(ctx context.Context, hash string) (dirItems []DirItem, err error)

	Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []DirItem, err error)

	FileCount(ctx context.Context) (int, error)
	DirCount(ctx context.Context) (int, error)
	DirItemCount(ctx context.Context) (int, error)
	BranchCount(ctx context.Context) (int, error)
}
