package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) Remove(ctx context.Context, branchName string, splitPath ...string) (dao.Commit, dao.Branch, error) {
	return fs.Db.RemoveDirItem(ctx, branchName, splitPath)
}
