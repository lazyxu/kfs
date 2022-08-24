package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) Checkout(ctx context.Context, branchName string) (bool, error) {
	return fs.Db.NewBranch(ctx, branchName)
}

func (fs *KFS) BranchInfo(ctx context.Context, branchName string) (branch dao.IBranch, err error) {
	return fs.Db.BranchInfo(ctx, branchName)
}
