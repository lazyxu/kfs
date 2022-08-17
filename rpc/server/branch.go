package server

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleBranchCheckout(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	// read
	branchName, err := rpcutil.ReadString(conn)
	if err != nil {
		return err
	}
	exist, err := kfsCore.Checkout(context.TODO(), branchName)
	if err != nil {
		return err
	}

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.BranchResp{Exist: exist})
	if err != nil {
		return err
	}
	return nil
}

func handleBranchInfo(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	// read
	branchName, err := rpcutil.ReadString(conn)
	if err != nil {
		return err
	}
	branch, err := kfsCore.BranchInfo(context.TODO(), branchName)
	if err != nil {
		return
	}

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.BranchInfoResp{
		Name:        branch.GetName(),
		Description: branch.GetDescription(),
		CommitId:    branch.GetCommitId(),
		Size:        branch.GetSize(),
		Count:       branch.GetCount(),
	})
	if err != nil {
		return err
	}
	return nil
}
