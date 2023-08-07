package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"github.com/pierrec/lz4"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func handleUploadDirItem(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// read
	var req pb.UploadReq
	err := rpcutil.ReadProto(conn, &req)
	if err != nil {
		return err
	}
	if req.Dir != nil {
		pbDirItems := req.Dir.DirItem
		fmt.Println(pbDirItems)
		dirItems := make([]dao.DirItem, len(pbDirItems))
		for i, dirItem := range pbDirItems {
			dirItems[i] = dao.DirItem{
				Hash:       dirItem.Hash,
				Name:       dirItem.Name,
				Mode:       dirItem.Mode,
				Size:       dirItem.Size,
				Count:      dirItem.Count,
				TotalCount: dirItem.TotalCount,
				CreateTime: dirItem.CreateTime,
				ModifyTime: dirItem.ModifyTime,
				ChangeTime: dirItem.ChangeTime,
				AccessTime: dirItem.AccessTime,
			}
		}
		dir, err := kfsCore.Db.WriteDir(context.TODO(), dirItems)
		if err != nil {
			return err
		}
		fmt.Println("UploadDir", dir)

		// write
		err = rpcutil.WriteOK(conn)
		if err != nil {
			return err
		}
		err = rpcutil.WriteProto(conn, &pb.UploadResp{Dir: &pb.DirResp{
			Hash:       dir.Hash(),
			Size:       dir.Size(),
			Count:      dir.Count(),
			TotalCount: dir.TotalCount(),
		},
		})
		return nil
	}
	root := req.Root
	dirItem := root.DirItem
	commit, branch, err := kfsCore.Db.UpsertDirItem(context.TODO(), root.BranchName, core.FormatPath(root.Path), dao.DirItem{
		Hash:       dirItem.Hash,
		Name:       dirItem.Name,
		Mode:       dirItem.Mode,
		Size:       dirItem.Size,
		Count:      dirItem.Count,
		TotalCount: dirItem.TotalCount,
		CreateTime: dirItem.CreateTime,
		ModifyTime: dirItem.ModifyTime,
		ChangeTime: dirItem.ChangeTime,
		AccessTime: dirItem.AccessTime,
	})
	if err != nil {
		return err
	}
	fmt.Println("Upload finish", root.Path)

	// write
	err = rpcutil.WriteOK(conn)
	if err != nil {
		return err
	}
	err = rpcutil.WriteProto(conn, &pb.UploadResp{
		Branch: &pb.BranchCommitResp{
			Hash:     commit.Hash,
			CommitId: commit.Id,
			Size:     branch.Size,
			Count:    branch.Count,
		},
	})
	return err
}

func handleUpload(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

	hashBytes := make([]byte, 256/8)
	err := binary.Read(conn, binary.LittleEndian, hashBytes)
	if err != nil {
		println(conn.RemoteAddr().String(), "hashBytes", err.Error())
		return rpcutil.UnexpectedIfError(err)
	}
	hash := hex.EncodeToString(hashBytes)
	println("hash", hash)

	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		println(conn.RemoteAddr().String(), "size", err.Error())
		return rpcutil.UnexpectedIfError(err)
	}
	println(conn.RemoteAddr().String(), "size", size)

	// 1. What if the hash is the same but the size is different?
	// 2. What if the hash and size are the same, but the file content is different?
	exist, err := kfsCore.S.Write(hash, func(f io.Writer, hasher io.Writer) (e error) {
		_, e = conn.Write([]byte{1}) // not exist
		if e != nil {
			return rpcutil.UnexpectedIfError(e)
		}

		encoder, e := rpcutil.ReadString(conn)
		println(conn.RemoteAddr().String(), "encoder", len(encoder), encoder)

		w := io.MultiWriter(f, hasher)
		if encoder == "lz4" {
			r := lz4.NewReader(conn)
			_, e = io.CopyN(w, r, size)
		} else {
			_, e = io.CopyN(w, conn, size)
		}
		println(conn.RemoteAddr().String(), "Copy")
		return rpcutil.UnexpectedIfError(e)
	})
	if err != nil {
		println(conn.RemoteAddr().String(), "Write", err.Error())
		return err
	}
	if exist {
		return nil
	}

	f := dao.NewFile(hash, uint64(size))
	err = kfsCore.Db.WriteFile(context.Background(), f)
	if err != nil {
		println(conn.RemoteAddr().String(), "Db.WriteFile", err.Error())
		return err
	}
	return nil
}
