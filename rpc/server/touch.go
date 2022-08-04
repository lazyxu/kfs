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
			rpcutil.WriteInvalid(conn, err)
		}
	}()

	// read
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
		CreateTime: req.Time,
		ModifyTime: req.Time,
		ChangeTime: req.Time,
		AccessTime: req.Time,
	})
	if err != nil {
		return
	}
	fmt.Println("Socket.Touch finish", req.String())

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return
	}
	err = rpcutil.WriteProto(conn, &pb.TouchResp{
		Hash:     commit.Hash,
		CommitId: commit.Id,
		Size:     branch.Size,
		Count:    branch.Count,
	})
	if err != nil {
		return
	}

	// exit
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return
	}
}