package client

import (
	"context"
	"github.com/lazyxu/kfs/dao"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) List(ctx context.Context, branchName string, filePath string, onLength func(int64) error, onDirItem func(item dao.IDirItem) error) error {
	var resp pb.DirItem
	err := ReqRespN(fs.SocketServerAddr, rpcutil.CommandList, &pb.PathReq{
		BranchName: branchName,
		Path:       filePath,
	}, &resp, onLength, func() error {
		return onDirItem(&resp)
	})
	if err != nil {
		return err
	}
	return nil
}
