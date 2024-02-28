package server

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pierrec/lz4"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleUploadStart(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadStartReq
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}

	now := time.Now().UnixNano()
	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.UploadStartResp{
		UploadTime: uint64(now),
	})
	if err != nil {
		return err
	}
	return nil
}

func handleUploadStartDir(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadStartDirReq
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "UploadStartDir", req.DriverId, "/"+strings.Join(req.DirPath, "/"))

	err = kfsCore.Db.UpsertDriverFile(context.TODO(), dao.DriverFile{
		DriverId:       req.DriverId,
		DirPath:        req.DirPath,
		Name:           req.Name,
		Hash:           req.Hash,
		Mode:           req.Mode,
		Size:           req.Size,
		CreateTime:     req.CreateTime,
		ModifyTime:     req.ModifyTime,
		ChangeTime:     req.ChangeTime,
		AccessTime:     req.AccessTime,
		UploadDeviceId: req.UploadDeviceId,
		UploadTime:     req.UploadTime,
	})

	l := len(req.UploadReqDirItemCheckV3)
	hashList := make([]string, l)
	// TODO: check exists.
	dirItemChecks := make([]dao.DirItemCheck, l)
	for i, c := range req.UploadReqDirItemCheckV3 {
		dirItemChecks[i] = dao.DirItemCheck{
			Name:       c.Name,
			Size:       c.Size,
			ModifyTime: c.ModifyTime,
		}
	}
	err = kfsCore.Db.CheckExists(context.TODO(), req.DriverId, req.DirPath, dirItemChecks, hashList)
	if err != nil {
		return err
	}
	if err != nil {
		println(conn.RemoteAddr().String(), "UpsertDriverFile", err.Error())
		return err
	}
	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.UploadStartDirResp{
		Hash: hashList,
	})
	if err != nil {
		return err
	}
	return nil
}

func handleUploadEndDir(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadEndDirReq
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "UploadDir", req.DriverId, "/"+strings.Join(req.DirPath, "/"))
	err = kfsCore.Db.SetLivpForMovAndHeicOrJpgInDirPath(context.TODO(), req.DriverId, req.DirPath)
	if err != nil {
		fmt.Println("Upload.SetLivpForMovAndHeicOrJpgInDriver error", err.Error())
		return err
	}
	return nil
}

func handleUploadV3File(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadFileV3
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "UploadFile", req.Hash, humanize.IBytes(req.Size))

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
	err = kfsCore.Db.InsertFile(context.TODO(), req.Hash, req.Size)
	if err != nil {
		println(conn.RemoteAddr().String(), "InsertFile", err.Error())
		return err
	}
	err = kfsCore.Db.UpsertDriverFile(context.TODO(), dao.DriverFile{
		DriverId:       req.DriverId,
		DirPath:        req.DirPath,
		Name:           req.Name,
		Hash:           req.Hash,
		Mode:           req.Mode,
		Size:           req.Size,
		CreateTime:     req.CreateTime,
		ModifyTime:     req.ModifyTime,
		ChangeTime:     req.ChangeTime,
		AccessTime:     req.AccessTime,
		UploadDeviceId: req.UploadDeviceId,
		UploadTime:     req.UploadTime,
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "UpsertDriverFile", err.Error())
		return err
	}
	err = AnalyzeIfNoFileType(context.TODO(), kfsCore, req.Hash)
	if err != nil {
		println(conn.RemoteAddr().String(), "AnalyzeIfNoFileType", req.Hash, err.Error())
		return err
	}
	err = PluginUnzipIfLivp(context.TODO(), kfsCore, req.Hash, req.Name)
	if err != nil {
		println(conn.RemoteAddr().String(), "PluginUnzipIfLivp", req.Hash, err.Error())
		return err
	}
	return nil
}
