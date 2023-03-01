package core

import (
	"context"

	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) Checkout(ctx context.Context, branchName string) (bool, error) {
	return fs.Db.NewBranch(ctx, branchName)
}

func (fs *KFS) ResetBranch(ctx context.Context, branchName string) error {
	return fs.Db.ResetBranch(ctx, branchName)
}

func (fs *KFS) DeleteBranch(ctx context.Context, branchName string) error {
	return fs.Db.DeleteBranch(ctx, branchName)
}

func (fs *KFS) BranchInfo(ctx context.Context, branchName string) (branch dao.IBranch, err error) {
	return fs.Db.BranchInfo(ctx, branchName)
}

func (fs *KFS) BranchList(ctx context.Context) ([]dao.IBranch, error) {
	return fs.Db.BranchList(ctx)
}

func (fs *KFS) BranchListCb(ctx context.Context, onLength func(int) error, onElement func(item dao.IBranch) error) error {
	list, err := fs.Db.BranchList(ctx)
	if err != nil {
		return err
	}
	if onLength != nil {
		err = onLength(len(list))
		if err != nil {
			return err
		}
	}
	if onElement != nil {
		for _, element := range list {
			err = onElement(element)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
