package core

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func Checkout(ctx context.Context, addr string, branchName string) (bool, error) {
	kfsCore, _, err := New(addr)
	if err != nil {
		return false, err
	}
	defer kfsCore.Close()
	return kfsCore.BranchNew(ctx, branchName)
}

func BranchInfo(ctx context.Context, addr string, branchName string) (sqlite.IBranch, error) {
	kfsCore, _, err := New(addr)
	if err != nil {
		return nil, err
	}
	defer kfsCore.Close()
	return kfsCore.BranchInfo(ctx, branchName)
}
