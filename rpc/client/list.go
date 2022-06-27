package client

import (
	"context"
	"io"
	"strconv"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) List(ctx context.Context, branchName string, filePath string, onLength func(int) error, onDirItem func(item sqlite.IDirItem) error) error {
	conn, c, err := getGRPCClient(fs)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := c.List(ctx, &pb.PathReq{
		BranchName: branchName,
		Path:       filePath,
	})
	if err != nil {
		return err
	}
	isFirst := true
	for {
		dirItem, err := client.Recv()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			return nil
		}
		if isFirst {
			md, err := client.Header()
			if err != nil {
				return err
			}
			length, err := strconv.Atoi(md.Get("length")[0])
			if err != nil {
				return err
			}
			err = onLength(length)
			if err != nil {
				return err
			}
			isFirst = false
		}
		err = onDirItem(dirItem)
		if err != nil {
			return err
		}
	}
}
