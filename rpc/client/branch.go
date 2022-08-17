package client

import (
	"context"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Checkout(ctx context.Context, branchName string) (bool, error) {
	var resp pb.BranchResp
	err := ReqStringResp(fs.SocketServerAddr, rpcutil.CommandBranchCheckout, branchName, &resp)
	if err != nil {
		return false, err
	}
	return resp.Exist, nil
}

func (fs *RpcFs) BranchInfo(ctx context.Context, branchName string) (sqlite.IBranch, error) {
	var resp pb.BranchInfoResp
	err := ReqStringResp(fs.SocketServerAddr, rpcutil.CommandBranchInfo, branchName, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
