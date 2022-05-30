package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func (fs *KFS) Remove(ctx context.Context, branchName string, splitPath ...string) (sqlite.Commit, sqlite.Branch, error) {
	return fs.Db.Remove(ctx, branchName, splitPath)
}
