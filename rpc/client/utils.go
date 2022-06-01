package client

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lazyxu/kfs/pb"
)

const fileChunkSize = 1024 * 1024

func withFS(fs GRPCFS, fn func(pb.KoalaFSClient) error) error {
	conn, err := grpc.Dial(fs.RemoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	return fn(c)
}

func withFS1[T any](fs GRPCFS, fn func(pb.KoalaFSClient) (T, error)) (T, error) {
	conn, err := grpc.Dial(fs.RemoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		var t T
		return t, err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	return fn(c)
}

func withFS2[T1 any, T2 any](fs GRPCFS, fn func(pb.KoalaFSClient) (T1, T2, error)) (T1, T2, error) {
	conn, err := grpc.Dial(fs.RemoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		var t1 T1
		var t2 T2
		return t1, t2, err
	}
	defer conn.Close()
	c := pb.NewKoalaFSClient(conn)
	return fn(c)
}

func SendContent(process core.UploadProcess, hash string, filename string, fn func(data []byte, isFirst bool, isLast bool) error) error {
	process.BeforeContent(hash, filename)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	isFirst := true
	for {
		chunk := make([]byte, 0, fileChunkSize)
		var n int64
		w := process.MultiWriter(bytes.NewBuffer(chunk))
		n, err = io.Copy(w, io.LimitReader(f, fileChunkSize))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = fn(chunk[:n], isFirst, n < fileChunkSize)
		isFirst = false
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

func FormatFilename(filename string) string {
	var name = []rune(filepath.Base(filename))
	if len(name) > 10 {
		name = append(name[:10], []rune("..")...)
	}
	return fmt.Sprintf("%-12s", string(name))
}
