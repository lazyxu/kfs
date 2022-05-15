package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lazyxu/kfs/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var uploadCmd = &cobra.Command{
	Use:     "upload",
	Example: "kfs-cli upload -b branchName -p path filePath",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runUpload,
}

func init() {
	uploadCmd.PersistentFlags().StringP(pathStr, "p", "", "")
}

const fileChunkSize = 1024 * 1024

func runUpload(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	remoteAddr := viper.GetString(remoteAddrStr)
	branchName := viper.GetString(branchNameStr)
	p := cmd.Flag(pathStr).Value.String()
	filename := args[0]
	fmt.Printf("remoteAddr=%s\n", remoteAddr)
	fmt.Printf("branch=%s\n", branchName)
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	ctx := context.Background()
	client, err := c.Upload(ctx)
	if err != nil {
		return
	}
	info, err := os.Stat(filename)
	if err != nil {
		return
	}
	hash, err := getFileHash(filename)
	if err != nil {
		return
	}
	now := uint64(time.Now().UnixNano())
	err = client.Send(&pb.UploadReq{Header: &pb.UploadReqHeader{
		BranchName: branchName,
		Path:       p,
		Hash:       hash,
		Mode:       uint64(info.Mode()),
		Size:       uint64(info.Size()),
		CreateTime: now,
		ModifyTime: uint64(info.ModTime().UnixNano()),
		ChangeTime: now,
		AccessTime: now,
	}})
	if err != nil {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	chunk := make([]byte, fileChunkSize)
	for {
		var n int
		n, err = f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		fmt.Printf("upload %d/%d\n", n, info.Size())
		err = client.Send(&pb.UploadReq{Bytes: chunk[:n]})
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if n < fileChunkSize {
			break
		}
	}
	resp, err := client.CloseAndRecv()
	if err == io.EOF {
		err = nil
	}
	if resp.Exist {
		fmt.Printf("the file already exists and does not need to be uploaded again\n")
	}
}

func getFileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
