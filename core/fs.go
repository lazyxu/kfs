package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type FS interface {
	Checkout(ctx context.Context, branchName string) (bool, error)
	BranchInfo(ctx context.Context, branchName string) (branch sqlite.IBranch, err error)
	List(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item sqlite.IDirItem) error) error
	Upload(ctx context.Context, branchName string, dstPath string, srcPath string, uploadProcess UploadProcess, concurrent int) (sqlite.Commit, sqlite.Branch, error)
	Reset(ctx context.Context) error
}
