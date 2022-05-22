package main

import (
	"fmt"
	"io"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
)

func (s *KoalaFSServer) Upload(server pb.KoalaFS_UploadServer) (err error) {
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return
	}
	defer kfsCore.Close()
	req := &pb.UploadReq{}
	req, err = server.Recv()
	if err != nil {
		return
	}
	h := req.Header
	fmt.Println("Upload", h)
	exist, commit, err := kfsCore.Upload(server.Context(), func(f io.Writer, hasher io.Writer) error {
		for {
			req, err = server.Recv()
			if req != nil {
				println("upload", req.IsLast, len(req.Bytes))
			}
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				return nil
			}
			_, err = hasher.Write(req.Bytes)
			if err != nil {
				return nil
			}
			_, err = f.Write(req.Bytes)
			if err != nil {
				return nil
			}
		}
	}, h.BranchName, h.Metadata.Path, h.Metadata.Hash,
		h.Metadata.Size, h.Metadata.Mode, h.Metadata.CreateTime, h.Metadata.ModifyTime, h.Metadata.ChangeTime, h.Metadata.AccessTime)
	if err != nil {
		return
	}
	err = server.SendAndClose(&pb.UploadResp{
		Exist:    exist,
		CommitId: commit.Id,
		Hash:     commit.Hash,
	})
	return
}