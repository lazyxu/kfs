package server

import (
	"context"
	"fmt"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) BranchCheckout(ctx context.Context, req *pb.BranchReq) (resp *pb.BranchResp, err error) {
	resp = new(pb.BranchResp)
	fmt.Println("BranchCheckout", req)
	resp.Exist, err = core.Checkout(ctx, s.kfsRoot, req.BranchName)
	if err != nil {
		return
	}
	return
}

func (s *KoalaFSServer) BranchInfo(ctx context.Context, req *pb.BranchInfoReq) (resp *pb.BranchInfoResp, err error) {
	resp = new(pb.BranchInfoResp)
	fmt.Println("BranchInfo", req)
	branch, err := core.BranchInfo(ctx, s.kfsRoot, req.BranchName)
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
