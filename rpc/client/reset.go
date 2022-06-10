package client

import (
	"context"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Reset(ctx context.Context, branchName string) error {
	return withFS(fs, func(c pb.KoalaFSClient) error {
		_, err := c.Reset(ctx, &pb.BranchReq{BranchName: branchName})
		return err
	})
}
