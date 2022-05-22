package list

import (
	"context"
	"io"
	"strconv"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func remote(ctx context.Context, addr string, branchName string, p string, isHumanize bool) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	client, err := c.List(ctx, &pb.PathReq{
		BranchName: branchName,
		Path:       p,
	})
	if err != nil {
		return err
	}
	isFirst := true
	for {
		dirItem := &pb.FileInfo{}
		dirItem, err = client.Recv()
		if err != nil && err != io.EOF {
			return err
		}
		isEOF := false
		if err == io.EOF {
			isEOF = true
			err = nil
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
			printHeader(length)
			isFirst = false
			printBody(dirItem, isHumanize)
		}
		if isEOF {
			return nil
		}
	}
}
