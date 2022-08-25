package core

import (
	"context"
	"os"
)

func (fs *KFS) Reset(ctx context.Context, branchName string) error {
	if fs.isSqlite {
		err := fs.Close()
		if err != nil {
			return err
		}
		err = os.RemoveAll(fs.root)
		if err != nil {
			return err
		}
		// TODO: fix reset with mysql
		kfs, _, err := NewWithSqlite(fs.root, fs.newStorage)
		if err != nil {
			return err
		}
		fs.Db = kfs.Db
		fs.S = kfs.S
	}
	_, err := fs.Checkout(ctx, branchName)
	return err
}
