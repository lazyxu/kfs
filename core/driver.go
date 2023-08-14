package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) InsertDriver(ctx context.Context, branchName string, description string) (bool, error) {
	return fs.Db.InsertDriver(ctx, branchName, description)
}

func (fs *KFS) DeleteDriver(ctx context.Context, branchName string) error {
	return fs.Db.DeleteDriver(ctx, branchName)
}

func (fs *KFS) ListDriver(ctx context.Context) ([]dao.IDriver, error) {
	return fs.Db.ListDriver(ctx)
}
