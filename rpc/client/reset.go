package client

import (
	"context"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Reset(ctx context.Context) error {
	return withFS(fs, func(c pb.KoalaFSClient) error {
		_, err := c.Reset(ctx, &pb.Void{})
		return err
	})
}
