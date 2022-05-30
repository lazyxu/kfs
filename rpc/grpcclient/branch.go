package grpcclient

import (
	"context"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Checkout(ctx context.Context, branchName string) (bool, error) {
	return withFS1[bool](fs, func(c pb.KoalaFSClient) (bool, error) {
		resp, err := c.BranchCheckout(ctx, &pb.BranchReq{
			BranchName: branchName,
		})
		return resp.Exist, err
	})
}

func (fs GRPCFS) BranchInfo(ctx context.Context, branchName string) (sqlite.IBranch, error) {
	return withFS1[sqlite.IBranch](fs, func(c pb.KoalaFSClient) (sqlite.IBranch, error) {
		return c.BranchInfo(ctx, &pb.BranchInfoReq{
			BranchName: branchName,
		})
	})
}
