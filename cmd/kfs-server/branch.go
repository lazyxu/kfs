package main

import (
	"context"
	"fmt"

	core "github.com/lazyxu/kfs/core/local"

	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) BranchCheckout(ctx context.Context, req *pb.BranchReq) (resp *pb.BranchResp, err error) {
	resp = new(pb.BranchResp)
	fmt.Println("BranchCheckout", req)
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return
	}
	defer kfsCore.Close()
	resp.Exist, err = kfsCore.BranchNew(ctx, req.BranchName, req.Description)
	if err != nil {
		return
	}
	return
}

func (s *KoalaFSServer) BranchInfo(ctx context.Context, req *pb.BranchInfoReq) (resp *pb.BranchInfoResp, err error) {
	resp = new(pb.BranchInfoResp)
	fmt.Println("BranchInfo", req)
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return
	}
	defer kfsCore.Close()
	branch, err := kfsCore.BranchInfo(ctx, req.BranchName)
	if err != nil {
		return
	}
	resp = &pb.BranchInfoResp{
		Name:        branch.Name,
		Description: branch.Description,
		CommitId:    branch.CommitId,
		Size:        branch.Size,
		Count:       branch.Count,
	}
	return
}
