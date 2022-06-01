package client

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/core"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Upload(ctx context.Context, branchName string, dstPath string, srcPath string, uploadProcess core.UploadProcess, concurrent bool) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	return withFS2[sqlite.Commit, sqlite.Branch](fs,
		func(c pb.KoalaFSClient) (commit sqlite.Commit, branch sqlite.Branch, err error) {
			client, err := c.Upload(ctx)
			if err != nil {
				return
			}
			srcPath, err = filepath.Abs(srcPath)
			if err != nil {
				return
			}

			walker := storage.NewWalker[sqlite.FileOrDir](ctx, srcPath, &uploadVisitor{
				client:        client,
				uploadProcess: uploadProcess,
				concurrent:    concurrent,
			})
			scanResp, err := walker.Walk(concurrent)
			if err != nil {
				return
			}
			info, err := os.Stat(srcPath)
			if err != nil {
				return
			}
			fileOrDir := scanResp.(sqlite.FileOrDir)
			modifyTime := uint64(info.ModTime().UnixNano())
			err = client.Send(&pb.UploadReq{
				Root: &pb.UploadReqRoot{
					BranchName: branchName,
					Path:       dstPath,
					DirItem: &pb.DirItem{
						Hash:       fileOrDir.Hash(),
						Name:       filepath.Base(dstPath),
						Mode:       uint64(info.Mode()),
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
			return sqlite.Commit{
					Id:   resp.Branch.CommitId,
					Hash: resp.Branch.Hash,
				}, sqlite.Branch{
					Name:     branchName,
					CommitId: resp.Branch.CommitId,
					Size:     resp.Branch.Size,
					Count:    resp.Branch.Count,
				}, nil
		})
}
