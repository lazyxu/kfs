package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func remote(ctx context.Context, addr string, branchName string, srcPath string, dstPath string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	client, err := c.Backup(ctx)
	if err != nil {
		return err
	}
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	backupCtx := storage.NewWalkerCtx[sqlite.FileOrDir](ctx, srcPath, &uploadVisitor{
		client:     client,
		branchName: branchName,
		backupPath: srcPath,
	})
	scanResp, err := backupCtx.Scan()
	if err != nil {
		return err
	}
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	fileOrDir := scanResp.(sqlite.FileOrDir)
	modifyTime := uint64(info.ModTime().UnixNano())
	err = client.Send(&pb.BackupReq{
		Root: &pb.BackupReqRoot{
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
		return err
	}
	resp, err := client.Recv()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp.Branch)
	return nil
}
