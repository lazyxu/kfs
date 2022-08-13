package server

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

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
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	fmt.Println("Socket.List", req.String())
	err = kfsCore.List(context.TODO(), req.BranchName, req.Path, func(n int) error {
		return binary.Write(conn, binary.LittleEndian, int64(n))
	}, func(dirItem sqlite.IDirItem) error {
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
