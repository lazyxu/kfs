package client

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Touch(ctx context.Context, branchName string, filePath string) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return
	}
	defer conn.Close()

	// TODO: upload empty file.

	fileOrDir := sqlite.NewFileByBytes(nil)
	modifyTime := uint64(time.Now().UnixNano())
	client, err := c.Upload(ctx)
	if err != nil {
		return
	}
	err = client.Send(&pb.UploadReq{
		Root: &pb.UploadReqRoot{
			BranchName: branchName,
			Path:       filePath,
			DirItem: &pb.DirItem{
				Hash:       fileOrDir.Hash(),
				Name:       filepath.Base(filePath),
				Mode:       uint64(os.FileMode(0o600)),
				Size:       fileOrDir.Size(),
				Count:      fileOrDir.Count(),
				TotalCount: fileOrDir.TotalCount(),
				CreateTime: modifyTime,
				ModifyTime: modifyTime,
				ChangeTime: modifyTime,
				AccessTime: modifyTime,
			},
		},
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return
	}
	resp, err := client.Recv()
	if err != nil {
		return
	}
	err = client.CloseSend()
	if err != nil {
		return
	}
	return sqlite.Commit{
			Id:   resp.Branch.CommitId,
			Hash: resp.Branch.Hash,
		}, sqlite.Branch{
			Name:     branchName,
			CommitId: resp.Branch.CommitId,
			Size:     resp.Branch.Size,
			Count:    resp.Branch.Count,
		}, nil
}
