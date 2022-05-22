package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func remote(ctx context.Context, addr string, filename string, branchName string, p string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	client, err := c.Upload(ctx)
	if err != nil {
		return err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	bar := NewProcessBar(info, filename)
	defer bar.Close()
	hash, err := SendHeader(bar, filename, info, p, func(metadata *pb.UploadReqMetadata) error {
		return client.Send(&pb.UploadReq{Header: &pb.UploadReqHeader{
			Metadata:   metadata,
			BranchName: branchName,
		}})
	})
	if err != nil {
		return err
	}
	err = SendContent(bar, hash, filename, func(data []byte, isLast bool) error {
		return client.Send(&pb.UploadReq{Bytes: data, IsLast: isLast})
	})
	if err != nil {
		return err
	}
	resp, err := client.CloseAndRecv()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return err
	}
	fmt.Println("branch updated with commit " + strconv.Itoa(int(resp.CommitId)) +
		" and hash " + resp.Hash[:4])
	return nil
}

func SendContent(bar *progressbar.ProgressBar, hash string, filename string, fn func(data []byte, isLast bool) error) error {
	bar.Reset()
	bar.Describe("[2/2][" + hash[0:4] + "] " + FormatFilename(filename))
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for {
		chunk := make([]byte, 0, fileChunkSize)
		var n int64
		w := io.MultiWriter(bytes.NewBuffer(chunk), bar)
		n, err = io.Copy(w, io.LimitReader(f, fileChunkSize))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = fn(chunk[:n], n < fileChunkSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n < fileChunkSize {
			break
		}
	}
	return nil
}

func NewProcessBar(info os.FileInfo, filename string) *progressbar.ProgressBar {
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
		progressbar.OptionSetDescription("[1/2][hash] "+FormatFilename(filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}

func FormatFilename(filename string) string {
	var name = []rune(path.Base(filename))
	if len(name) > 10 {
		name = append(name[:10], []rune("..")...)
	}
	return fmt.Sprintf("%-12s", string(name))
}
