package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) ListCb(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item dao.IDirItem) error) error {
	dirItems, err := fs.Db.List(ctx, branchName, FormatPath(filePath))
	if err != nil {
		return err
	}
	if onLength != nil {
		err = onLength(len(dirItems))
		if err != nil {
			return err
		}
	}
	if onDirItem != nil {
		for _, dirItem := range dirItems {
			err = onDirItem(dirItem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (fs *KFS) List(ctx context.Context, branchName string, filePath string) ([]dao.DirItem, error) {
	return fs.Db.List(ctx, branchName, FormatPath(filePath))
}

func (fs *KFS) ListDriverFile(ctx context.Context, driverId uint64, filePath []string) ([]dao.DriverFile, error) {
	return fs.Db.ListDriverFile(ctx, driverId, filePath)
}
