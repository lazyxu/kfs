package client

import (
	"context"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func (fs *RpcFs) Reset(ctx context.Context, branchName string) error {
	return ReqString(fs.SocketServerAddr, rpcutil.CommandReset, branchName)
}
