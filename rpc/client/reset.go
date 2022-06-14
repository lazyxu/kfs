package client

import (
	"context"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Reset(ctx context.Context, branchName string) error {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = c.Reset(ctx, &pb.BranchReq{BranchName: branchName})
	return err
}
