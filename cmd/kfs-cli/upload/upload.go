package upload

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/lazyxu/kfs/cmd/kfs-cli/utils"

	"github.com/schollz/progressbar/v3"

	"github.com/lazyxu/kfs/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Cmd = &cobra.Command{
	Use:     "upload",
	Example: "kfs-cli upload -b branchName -p path filePath",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runUpload,
}

func init() {
	Cmd.PersistentFlags().StringP(utils.PathStr, "p", "", "")
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
	remoteAddr := viper.GetString(utils.ServerAddrStr)
	branchName := viper.GetString(utils.BranchNameStr)
	p := cmd.Flag(utils.PathStr).Value.String()
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
	bar := progressbar.NewOptions(int(info.Size()),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(20*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetDescription("[1/2][hash] "+formatFilename(filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	hash, err := getFileHash(bar, filename)
	if err != nil {
		return
	}
	modifyTime := uint64(info.ModTime().UnixNano())
	err = client.Send(&pb.UploadReq{Header: &pb.UploadReqHeader{
		BranchName: branchName,
		Path:       p,
		Hash:       hash,
		Mode:       uint64(info.Mode()),
		Size:       uint64(info.Size()),
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	}})
	if err != nil {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	bar.Reset()
	bar.Describe("[2/2][" + hash[0:4] + "] " + formatFilename(filename))
	chunk := make([]byte, 0, fileChunkSize)
	for {
		var n int64
		w := io.MultiWriter(bytes.NewBuffer(chunk), bar)
		n, err = io.Copy(w, io.LimitReader(f, fileChunkSize))
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
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
	_ = bar.Close()
	resp, err := client.CloseAndRecv()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return
	}
	fmt.Println("branch updated with commit " + strconv.Itoa(int(resp.CommitId)) +
		" and hash " + resp.Hash[:4])
}

func formatFilename(filename string) string {
	var name = []rune(path.Base(filename))
	if len(name) > 10 {
		name = append(name[:10], []rune("..")...)
	}
	return fmt.Sprintf("%-12s", string(name))
}

func getFileHash(bar io.Writer, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := sha256.New()
	w := io.MultiWriter(hash, bar)
	_, err = io.Copy(w, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
