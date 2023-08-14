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

	InsertDriver(ctx context.Context, driverName string, description string) (exist bool, err error)
	DeleteDriver(ctx context.Context, driverName string) error
	ListDriver(ctx context.Context) (drivers []IDriver, err error)

	InsertFile(ctx context.Context, hash string, size uint64) error

	UpsertDriverFile(ctx context.Context, f DriverFile) error
	ListDriverFile(ctx context.Context, driverName string, filePath []string) (files []DriverFile, err error)
	GetDriverFile(ctx context.Context, driverName string, splitPath []string) (file DriverFile, err error)

	InsertNullExif(ctx context.Context, hash string) (exist bool, err error)
	InsertExif(ctx context.Context, hash string, e ExifData) (exist bool, err error)
	ListExpectExif(ctx context.Context) (hashList []string, err error)
	ListExpectExifCb(ctx context.Context, cb func(hash string)) (err error)
}

func DatabaseNewFunc(dataSourceName string, newDB func(dataSourceName string) (Database, error)) func() (Database, error) {
	return func() (Database, error) {
		return newDB(dataSourceName)
	}
}
