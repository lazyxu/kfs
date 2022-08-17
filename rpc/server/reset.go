package server

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

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
