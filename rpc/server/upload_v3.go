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
	"github.com/lazyxu/kfs/db/sqlc"
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
	ctx := context.TODO()
	// read
	var req pb.UploadStartDirReq
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "UploadStartDir", req.ParentID, "/"+strings.Join(req.DirPath, "/"))

	var userID int64 = 1
	var parentID dao.NullInt64 = sqlc.NInt64Null()
	if !req.IsRoot {
		parentID = sqlc.NInt64(req.ParentID)
		_, err = kfsCore.Q.SetLogicFileByParentIdPath(ctx, userID, parentID,
			req.DirPath, req.Name, sqlc.NInt64Null(), req.Mode, req.Size,
			time.Unix(0, req.CreateTime), time.Unix(0, req.ModifyTime),
			time.Unix(0, req.ChangeTime), time.Unix(0, req.AccessTime), false)
		// err = kfsCore.Q.SetLogicFileByParentIdPath(context.TODO(), dao.DriverFile{
		// 	DriverId:       req.FileID,
		// 	DirPath:        req.DirPath,
		// 	Name:           req.Name,
		// 	Hash:           req.Hash,
		// 	Mode:           req.Mode,
		// 	Size:           req.Size,
		// 	CreateTime:     req.CreateTime,
		// 	ModifyTime:     req.ModifyTime,
		// 	ChangeTime:     req.ChangeTime,
		// 	AccessTime:     req.AccessTime,
		// 	UploadDeviceId: req.UploadDeviceId,
		// 	UploadTime:     req.UploadTime,
		// }, false)
		if err != nil {
			println(conn.RemoteAddr().String(), "UpsertDriverFile", err.Error())
			return err
		}
	}

	// l := len(req.UploadReqDirItemCheckV3)
	// hashList := make([]string, l)
	// // TODO: check exists.
	// dirItemChecks := make([]dao.DirItemCheck, l)
	// for i, c := range req.UploadReqDirItemCheckV3 {
	// 	dirItemChecks[i] = dao.DirItemCheck{
	// 		Name:       c.Name,
	// 		Size:       c.Size,
	// 		ModifyTime: c.ModifyTime,
	// 	}
	// }
	lfs, err := kfsCore.Q.ListLogicFilesByParentId(ctx, sqlc.ListLogicFilesByParentIdParams{
		UserID:   userID,
		ParentID: parentID,
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "ListLogicFilesByParentId", err.Error())
		return err
	}
	notExistInDir := make([]int32, 0, len(req.UploadReqDirItemCheckV3))
	// TODO: 当前只添加覆盖文件，不会删除文件
	type LogicFileInfo struct {
		Size       int64
		ModifyTime time.Time
	}
	lfi := make(map[string]LogicFileInfo, len(lfs))
	for _, lf := range lfs {
		lfi[lf.Name] = LogicFileInfo{
			Size:       lf.Size,
			ModifyTime: lf.ModifyTime,
		}
	}
	for i, item := range req.UploadReqDirItemCheckV3 {
		exist := false
		if info, ok := lfi[item.Name]; ok {
			if info.Size == item.Size && info.ModifyTime.Equal(time.Unix(0, item.ModifyTime)) {
				exist = true
			}
		}
		if !exist {
			notExistInDir = append(notExistInDir, int32(i))
		}
	}

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.UploadStartDirResp{
		IndexForNotExistInDir: notExistInDir,
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
	println(conn.RemoteAddr().String(), "UploadDir", req.ParentID, "/"+strings.Join(req.DirPath, "/"))
	err = kfsCore.Q.SetLivpForMovAndHeicOrJpgByParentId(context.TODO(), sqlc.SetLivpForMovAndHeicOrJpgByParentIdParams{
		ParentID: sqlc.NInt64(req.ParentID),
	})
	if err != nil {
		fmt.Println("Upload.SetLivpForMovAndHeicOrJpgInDriver error", err.Error())
		return err
	}
	return nil
}

func handleUploadV3File(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	ctx := context.TODO()
	// read
	var req pb.UploadFileV3
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	println(conn.RemoteAddr().String(), "UploadFile", req.MD5, req.SHA256, humanize.IBytes(uint64(req.Size)))

	_, err = kfsCore.S2.AddReader(req.Size, req.MD5, req.SHA256, func() (io.Reader, error) {
		_, err := conn.Write([]byte{1}) // not exist
		if err != nil {
			return nil, rpcutil.UnexpectedIfError(err)
		}

		encoder, err := rpcutil.ReadString(conn)
		if err != nil {
			return nil, rpcutil.UnexpectedIfError(err)
		}
		println(conn.RemoteAddr().String(), "encoder", len(encoder), encoder)

		var r io.Reader = conn
		if encoder == "lz4" {
			r = lz4.NewReader(conn)
		}
		println(conn.RemoteAddr().String(), "CopyStart", req.Size)
		return r, nil
	})

	if err != nil {
		println(conn.RemoteAddr().String(), "Write", err.Error())
		return err
	}
	res, err := kfsCore.Q.UpsertFile(ctx, sqlc.UpsertFileParams{
		Size:   req.Size,
		Md5:    req.MD5,
		Sha256: req.SHA256,
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "UpsertFile", err.Error())
		return err
	}
	fileID, err := res.LastInsertId()
	if err != nil {
		println(conn.RemoteAddr().String(), "LastInsertId", err.Error())
		return err
	}
	var userID int64 = 1
	var parentID dao.NullInt64 = sqlc.NInt64Null()
	if req.ParentID != 0 {
		parentID = sqlc.NInt64(req.ParentID)
	}
	_, err = kfsCore.Q.SetLogicFileByParentIdPath(ctx, userID, parentID,
		req.DirPath, req.Name, sqlc.NInt64(fileID), req.Mode, req.Size,
		time.Unix(0, req.CreateTime), time.Unix(0, req.ModifyTime),
		time.Unix(0, req.ChangeTime), time.Unix(0, req.AccessTime), false)
	if err != nil {
		println(conn.RemoteAddr().String(), "UpsertDriverFile", err.Error())
		return err
	}
	f := sqlc.File{
		ID: fileID, Size: req.Size, Md5: req.MD5, Sha256: req.SHA256,
	}
	err = AnalyzeMetadataNotAnalyzed(context.TODO(), kfsCore, f)
	if err != nil {
		println(conn.RemoteAddr().String(), "AnalyzeIfNoFileType", fileID, err.Error())
		return err
	}
	err = PluginAnalyzeLivpNotAnalyzed(context.TODO(), kfsCore, req.Name, f)
	if err != nil {
		println(conn.RemoteAddr().String(), "PluginUnzipIfLivp", fileID, err.Error())
		return err
	}
	return nil
}
