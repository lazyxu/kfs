package server

import (
	"context"
	"fmt"
	"github.com/pierrec/lz4"
	"io"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleUploadV2Dir(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadReqV2
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	err = kfsCore.Db.UpsertDriverFile(context.TODO(), dao.DriverFile{
		DriverName: req.DriverName,
		FilePath:   core.FormatPathV2(req.FilePath),
		Version:    0,
		Hash:       req.Hash,
		Mode:       req.Mode,
		Size:       req.Size,
		CreateTime: req.CreateTime,
		ModifyTime: req.ModifyTime,
		ChangeTime: req.ChangeTime,
		AccessTime: req.AccessTime,
	})
	if err != nil {
		return err
	}
	fmt.Println("Upload finish", req.FilePath)

	return nil
}

func handleUploadV2File(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadReqV2
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "hash", len(req.Hash), req.Hash)

	// 1. What if the hash is the same but the size is different?
	// 2. What if the hash and size are the same, but the file content is different?
	_, err = kfsCore.S.Write(req.Hash, func(f io.Writer, hasher io.Writer) (e error) {
		_, e = conn.Write([]byte{1}) // not exist
		if e != nil {
			return rpcutil.UnexpectedIfError(e)
		}

		encoder, e := rpcutil.ReadString(conn)
		println(conn.RemoteAddr().String(), "encoder", len(encoder), encoder)

		println(conn.RemoteAddr().String(), "CopyStart", req.Size)
		w := io.MultiWriter(f, hasher)
		var n int64
		if encoder == "lz4" {
			r := lz4.NewReader(conn)
			n, e = io.CopyN(w, r, int64(req.Size))
		} else {
			n, e = io.CopyN(w, conn, int64(req.Size))
		}
		println(conn.RemoteAddr().String(), "CopyEnd", n)
		return rpcutil.UnexpectedIfError(e)
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "Write", err.Error())
		return err
	}
	err = kfsCore.Db.UpsertDriverFile(context.TODO(), dao.DriverFile{
		DriverName: req.DriverName,
		FilePath:   core.FormatPathV2(req.FilePath),
		Version:    0,
		Hash:       req.Hash,
		Mode:       req.Mode,
		Size:       req.Size,
		CreateTime: req.CreateTime,
		ModifyTime: req.ModifyTime,
		ChangeTime: req.ChangeTime,
		AccessTime: req.AccessTime,
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "UpsertDriverFile", err.Error())
		return err
	}

	return nil
}
