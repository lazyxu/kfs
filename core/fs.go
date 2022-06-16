package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type UploadConfig struct {
	Encoder       string
	UploadProcess UploadProcess
	Concurrent    int
	Verbose       bool
}

type FS interface {
	Checkout(ctx context.Context, branchName string) (bool, error)
	BranchInfo(ctx context.Context, branchName string) (branch sqlite.IBranch, err error)
	List(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item sqlite.IDirItem) error) error
	Upload(ctx context.Context, branchName string, dstPath string, srcPath string, config UploadConfig) (sqlite.Commit, sqlite.Branch, error)
	Reset(ctx context.Context, branchName string) error
	Close() error
	Download(ctx context.Context, branchName string, dstPath string, srcPath string, config UploadConfig) (string, error)
}
