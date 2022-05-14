package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/lazyxu/kfs/pb"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"github.com/dustin/go-humanize"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Short:   "ls list files",
	Example: "kfs-cli ls .",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runList,
}

const (
	kfsRootStr    = "kfs-root"
	backupPathStr = "backup-path"
	branchNameStr = "branch-name"
	pathStr       = "path"
)

func runList(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	remoteAddr := viper.GetString(remoteAddrStr)
	branchName := viper.GetString(branchNameStr)
	fmt.Printf("remoteAddr=%s\n", remoteAddr)
	fmt.Printf("branch=%s\n", branchName)
	p := ""
	if len(args) != 0 {
		p = args[0]
	}
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	ctx := context.Background()
	client, err := c.List(ctx, &pb.PathReq{
		BranchName: branchName,
		Path:       p,
	})
	if err != nil {
		return
	}
	isFirst := true
	for {
		dirItem := &pb.FileInfo{}
		dirItem, err = client.Recv()
		if err != nil && err != io.EOF {
			return
		}
		isEOF := false
		if err == io.EOF {
			isEOF = true
		}
		if isFirst {
			var md metadata.MD
			md, err = client.Header()
			if err != nil {
				return
			}
			length, err := strconv.Atoi(md.Get("length")[0])
			if err != nil {
				return
			}
			fmt.Printf("total %d\n", length)
			if length != 0 {
				fmt.Printf("mode      \tcount\ttotalCount\thash\tsize\tmodifyTime         \tname\n")
			}
			isFirst = false
		}
		if isEOF {
			return
		}
		modifyTime := time.Unix(0, int64(dirItem.ModifyTime)).Format("2006-01-02 15:04:05")
		fmt.Printf("%s\t%s\t     %s\t%s\t%s\t%s\t%s\n",
			os.FileMode(dirItem.Mode).String(),
			formatCount(dirItem.Mode, dirItem.Count), formatCount(dirItem.Mode, dirItem.TotalCount), dirItem.Hash[:4],
			humanize.Bytes(dirItem.Size), modifyTime, dirItem.Name)
	}
}

func formatCount(mode uint64, count uint64) string {
	if !os.FileMode(mode).IsDir() {
		return strings.Repeat(" ", 5)
	}
	return fmt.Sprintf("%5d", count)
}
