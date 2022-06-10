package server

import (
	"context"

	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Reset(ctx context.Context, req *pb.BranchReq) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	err = s.kfsCore.Reset(ctx, req.BranchName)
	return
}
