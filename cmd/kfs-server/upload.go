package main

import (
	"fmt"
	"io"
	"path/filepath"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Upload(server pb.KoalaFS_UploadServer) (err error) {
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return
	}
	defer kfsCore.Close()
	req, err := server.Recv()
	if err != nil {
		return
	}
	exist, err := kfsCore.S.WriteFn(req.File.Hash, func(f io.Writer, hasher io.Writer) error {
		for {
			file := req.File
			println("Upload", len(file.Hash), len(file.Bytes), file.IsLastChunk)
			_, err = hasher.Write(file.Bytes)
			if err != nil {
				return err
			}
			_, err = f.Write(file.Bytes)
			if err != nil {
				return err
			}
			if file.IsLastChunk {
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
	for req.File != nil { // skip if file exists
		req, err = server.Recv()
		if err != nil {
			return err
		}
	}
	root := req.Root
	dirItem := root.DirItem
	fmt.Println("Upload", req)
	ext := filepath.Ext(dirItem.Name)
	f := sqlite.NewFile(dirItem.Hash, dirItem.Size, ext)
	err = kfsCore.Db.WriteFile(server.Context(), f)
	if err != nil {
		return err
	}
	commit, branch, err := kfsCore.Db.UpsertDirItem(server.Context(), root.BranchName, core.FormatPath(root.Path), sqlite.DirItem{
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
	err = server.SendAndClose(&pb.UploadResp{
		Exist: exist,
		Branch: &pb.BranchCommitResp{
			Hash:     commit.Hash,
			CommitId: commit.Id,
			Size:     branch.Size,
			Count:    branch.Count,
		},
	})
	return
}
