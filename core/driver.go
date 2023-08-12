package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) NewDriver(ctx context.Context, branchName string, description string) (bool, error) {
	return fs.Db.NewDriver(ctx, branchName, description)
}

func (fs *KFS) DeleteDriver(ctx context.Context, branchName string) error {
	return fs.Db.DeleteDriver(ctx, branchName)
}

func (fs *KFS) DriverList(ctx context.Context) ([]dao.IDriver, error) {
	return fs.Db.DriverList(ctx)
}
