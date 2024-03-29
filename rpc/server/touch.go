package server

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/lazyxu/kfs/dao"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleTouch(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	// read
	var req pb.TouchReq
	err = rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	fileOrDir := dao.NewFileByBytes(nil)
	_, err = kfsCore.S.Write(fileOrDir.Hash(), func(f io.Writer, hasher io.Writer) error {
		return nil
	})
	if err != nil {
		return err
	}
	commit, branch, err := kfsCore.Db.UpsertDirItem(context.TODO(), req.BranchName, core.FormatPath(req.Path), dao.DirItem{
		Hash:       fileOrDir.Hash(),
		Name:       filepath.Base(req.Path),
		Mode:       req.Mode,
		Size:       fileOrDir.Size(),
		Count:      fileOrDir.Count(),
		TotalCount: fileOrDir.TotalCount(),
		CreateTime: req.Time,
		ModifyTime: req.Time,
		ChangeTime: req.Time,
		AccessTime: req.Time,
	})
	if err != nil {
		return err
	}
	fmt.Println("Socket.Touch finish", req.String())

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.TouchResp{
		Hash:     commit.Hash,
		CommitId: commit.Id,
		Size:     branch.Size,
		Count:    branch.Count,
	})
	if err != nil {
		return err
	}
	return nil
}
