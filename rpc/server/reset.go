package server

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func (s *KoalaFSServer) Reset(ctx context.Context, req *pb.BranchReq) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	err = s.kfsCore.Reset(ctx, req.BranchName)
	return
}

func handleReset(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	branchName, err := rpcutil.ReadString(conn)
	if err != nil {
		return err
	}
	err = kfsCore.Reset(context.TODO(), branchName)
	if err != nil {
		println(conn.RemoteAddr().String(), "Reset", err.Error())
		return err
	}
	return nil
}
