package list

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func remote(addr string, branchName string, p string, human string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	ctx := context.Background()
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
		}
		if isEOF {
			return nil
		}
		modifyTime := time.Unix(0, int64(dirItem.ModifyTime)).Format("2006-01-02 15:04:05")
		if human == "true" {
			fmt.Printf("%s\t%s\t     %s\t%s\t%s\t%s\t%s\n",
				os.FileMode(dirItem.Mode).String(),
				formatCount(dirItem.Mode, dirItem.Count), formatCount(dirItem.Mode, dirItem.TotalCount), dirItem.Hash[:4],
				humanize.Bytes(dirItem.Size), modifyTime, dirItem.Name)
		} else {
			fmt.Printf("%s\t%s\t     %s\t%s\t%d\t%s\t%s\n",
				os.FileMode(dirItem.Mode).String(),
				formatCount(dirItem.Mode, dirItem.Count), formatCount(dirItem.Mode, dirItem.TotalCount), dirItem.Hash[:4],
				dirItem.Size, modifyTime, dirItem.Name)
		}
	}
}

func formatCount(mode uint64, count uint64) string {
	if !os.FileMode(mode).IsDir() {
		return strings.Repeat(" ", 5)
	}
	return fmt.Sprintf("%5d", count)
}
