package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func handleTouch(kfsCore *core.KFS, conn net.Conn) {
	var err error
	defer func() {
		if err != nil {
			rpcutil.WriteErrorExit(conn, err)
		}
	}()
	var req pb.TouchReq
	err = rpcutil.ReadProto(conn, &req)
	if err != nil {
		return
	}
	fileOrDir := sqlite.NewFileByBytes(nil)
	_, err = kfsCore.S.WriteFn(fileOrDir.Hash(), func(f io.Writer, hasher io.Writer) error {
		return nil
	})
	if err != nil {
		return
	}
	commit, branch, err := kfsCore.Db.UpsertDirItem(context.TODO(), req.BranchName, core.FormatPath(req.Path), sqlite.DirItem{
		Hash:       fileOrDir.Hash(),
		Name:       filepath.Base(req.Path),
		Mode:       req.Mode,
		Size:       fileOrDir.Size(),
		Count:      fileOrDir.Count(),
		TotalCount: fileOrDir.TotalCount(),
		CreateTime: req.CreateTime,
		ModifyTime: req.ModifyTime,
		ChangeTime: req.ChangeTime,
		AccessTime: req.AccessTime,
	})
	if err != nil {
		return
	}
	fmt.Println("Touch finish", req.Path)
	err = rpcutil.WriteProto(conn, &pb.TouchResp{
		Hash:     commit.Hash,
		CommitId: commit.Id,
		Size:     branch.Size,
		Count:    branch.Count,
	})
	if err != nil {
		return
	}
}
