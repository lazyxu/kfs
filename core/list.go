package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func (fs *KFS) List(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item sqlite.IDirItem) error) error {
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
