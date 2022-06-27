package client

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Checkout(ctx context.Context, branchName string) (bool, error) {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	resp, err := c.BranchCheckout(ctx, &pb.BranchReq{
		BranchName: branchName,
	})
	if err != nil {
		return false, err
	}
	return resp.Exist, nil
}

func (fs *RpcFs) BranchInfo(ctx context.Context, branchName string) (sqlite.IBranch, error) {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return c.BranchInfo(ctx, &pb.BranchInfoReq{
		BranchName: branchName,
	})
}
