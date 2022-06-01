package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func (fs *KFS) Checkout(ctx context.Context, branchName string) (bool, error) {
	return fs.Db.NewBranch(ctx, branchName)
}

func (fs *KFS) BranchInfo(ctx context.Context, branchName string) (branch sqlite.IBranch, err error) {
	return fs.Db.BranchInfo(ctx, branchName)
}
