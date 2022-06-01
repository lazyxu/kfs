package server

import (
	"context"
	"fmt"

	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) BranchCheckout(ctx context.Context, req *pb.BranchReq) (resp *pb.BranchResp, err error) {
	resp = new(pb.BranchResp)
	fmt.Println("BranchCheckout", req)
	resp.Exist, err = s.kfsCore.Checkout(ctx, req.BranchName)
	if err != nil {
		return
	}
	return
}

func (s *KoalaFSServer) BranchInfo(ctx context.Context, req *pb.BranchInfoReq) (resp *pb.BranchInfoResp, err error) {
	resp = new(pb.BranchInfoResp)
	fmt.Println("BranchInfo", req)
	branch, err := s.kfsCore.BranchInfo(ctx, req.BranchName)
	if err != nil {
		return
	}
	resp = &pb.BranchInfoResp{
		Name:        branch.GetName(),
		Description: branch.GetDescription(),
		CommitId:    branch.GetCommitId(),
		Size:        branch.GetSize(),
		Count:       branch.GetCount(),
	}
	return
}
