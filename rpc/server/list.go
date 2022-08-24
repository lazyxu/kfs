package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/lazyxu/kfs/dao"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/pb"
)

func handleList(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	// read
	var req pb.PathReq
	err = rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}

	// write
	fmt.Println("Socket.List", req.String())
	err = kfsCore.List(context.TODO(), req.BranchName, req.Path, func(n int) error {
		err = rpcutil.WriteOK(conn)
		if err != nil {
			return err
		}
		return binary.Write(conn, binary.LittleEndian, int64(n))
	}, func(dirItem dao.IDirItem) error {
		return rpcutil.WriteProto(conn, &pb.DirItem{
			Hash:       dirItem.GetHash(),
			Name:       dirItem.GetName(),
			Mode:       dirItem.GetMode(),
			Size:       dirItem.GetSize(),
			Count:      dirItem.GetCount(),
			TotalCount: dirItem.GetTotalCount(),
			CreateTime: dirItem.GetCreateTime(),
			ModifyTime: dirItem.GetModifyTime(),
			ChangeTime: dirItem.GetChangeTime(),
			AccessTime: dirItem.GetAccessTime(),
		})
	})
	if err != nil {
		return err
	}
	return nil
}
