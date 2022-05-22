package backup

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func remote(ctx context.Context, addr string, branchName string, dstPath string, backupPath string) error {
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
	backupPath, err = filepath.Abs(backupPath)
	if err != nil {
		return err
	}

	backupCtx := storage.NewWalkerCtx[sqlite.FileOrDir](ctx, backupPath, &uploadVisitor{
		client:     client,
		branchName: branchName,
		backupPath: backupPath,
	})
	_, err = backupCtx.Scan()
	if err != nil {
		return err
	}
	err = client.Send(&pb.BackupReq{
		Done: true,
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
	fmt.Println("branch updated with commit " + strconv.Itoa(int(resp.UploadResp.CommitId)) +
		" and hash " + resp.UploadResp.Hash[:4])
	return nil
}
