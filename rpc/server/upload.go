package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"github.com/pierrec/lz4"
	"io"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Upload(server pb.KoalaFS_UploadServer) (err error) {
	req := &pb.UploadReq{}
	var exist bool
	for {
		req, err = server.Recv()
		if err != nil {
			return err
		}
		if req.Root != nil {
			break
		}
		if req.File != nil {
			if req.File.Hash == "" {
				continue // file already exists, ignored
			}
			firstFileChunk := req.File
			exist, err = s.kfsCore.S.WriteFn(firstFileChunk.Hash, func(f io.Writer, hasher io.Writer) error {
				for {
					_, err = hasher.Write(req.File.Bytes)
					if err != nil {
						return err
					}
					_, err = f.Write(req.File.Bytes)
					if err != nil {
						return err
					}
					if req.File.IsLastChunk {
						return nil
					}
					req, err = server.Recv()
					if err != nil {
						return err
					}
				}
			})
			if err != nil {
				return
			}
			f := sqlite.NewFile(firstFileChunk.Hash, firstFileChunk.Size)
			err = s.kfsCore.Db.WriteFile(server.Context(), f)
			if err != nil {
				return
			}
			fmt.Println("Upload", f, exist)
			err = server.Send(&pb.UploadResp{Exist: exist})
			if err != nil {
				return
			}
		} else {
			pbDirItems := req.Dir.DirItem
			fmt.Println(pbDirItems)
			dirItems := make([]sqlite.DirItem, len(pbDirItems))
			for i, dirItem := range pbDirItems {
				dirItems[i] = sqlite.DirItem{
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
			var dir sqlite.Dir
			dir, err = s.kfsCore.Db.WriteDir(server.Context(), dirItems)
			fmt.Println("UploadDir", dir)
			err = server.Send(&pb.UploadResp{Dir: &pb.DirResp{
				Hash:       dir.Hash(),
				Size:       dir.Size(),
				Count:      dir.Count(),
				TotalCount: dir.TotalCount(),
			},
			})
		}
	}
	root := req.Root
	dirItem := root.DirItem
	commit, branch, err := s.kfsCore.Db.UpsertDirItem(server.Context(), root.BranchName, core.FormatPath(root.Path), sqlite.DirItem{
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
		return
	}
	fmt.Println("Upload finish", root.Path)
	err = server.Send(&pb.UploadResp{
		Branch: &pb.BranchCommitResp{
			Hash:     commit.Hash,
			CommitId: commit.Id,
			Size:     branch.Size,
			Count:    branch.Count,
		},
	})
	return
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

	exist, err := kfsCore.S.WriteFn(hash, func(f io.Writer, hasher io.Writer) (e error) {
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
		println(conn.RemoteAddr().String(), "WriteFn", err.Error())
		return err
	}
	if exist {
		return nil
	}

	f := sqlite.NewFile(hash, uint64(size))
	err = kfsCore.Db.WriteFile(context.Background(), f)
	if err != nil {
		println(conn.RemoteAddr().String(), "Db.WriteFile", err.Error())
		return err
	}
	return nil
}
