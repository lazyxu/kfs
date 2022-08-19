package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"io"
)

func handleOpen(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	// read
	var req pb.PathReq
	err = rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}

	// write
	fmt.Println("Socket.Open", req.String())
	mode, rc, dirItems, err := kfsCore.Open(context.TODO(), req.BranchName, req.Path)
	if err != nil {
		return err
	}
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, int64(mode))
	if err != nil {
		return err
	}
	if mode.IsRegular() {
		err = binary.Write(conn, binary.LittleEndian, rc.Size())
		if err != nil {
			return err
		}
		_, err = io.Copy(conn, rc)
		if err != nil {
			return err
		}
		return
	}
	err = binary.Write(conn, binary.LittleEndian, int64(len(dirItems)))
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		err = rpcutil.WriteProto(conn, &pb.DirItem{
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
		if err != nil {
			return err
		}
	}
	return nil
}
