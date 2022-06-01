package server

import (
	"fmt"
	"io"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Upload(server pb.KoalaFS_UploadServer) (err error) {
	req := &pb.UploadReq{}
	var exist bool
	fmt.Println("-----------")
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
			// TODO: upload dir
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
			fmt.Println("Upload", dir)
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
