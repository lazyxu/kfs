package localfs

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/lazyxu/kfs/warpper/grpcweb/rootdirectory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/proto"
	"github.com/lazyxu/kfs/warpper/grpcweb/pb"
	"github.com/sirupsen/logrus"
)

type LenWriter struct {
	n int64
}

func (w *LenWriter) Write(p []byte) (n int, err error) {
	w.n += int64(len(p))
	return len(p), nil
}

func getConn(f func(conn net.Conn) error) {

}

func (c *BackUpCtx) uploadFile(filePath string, size int64) (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:9877")
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return "", err
	}
	//h := sha256.New()
	//_, err = io.Copy(h, f)
	//if err != nil {
	//	return "", err
	//}
	//hash := h.Sum(nil)
	//logrus.Infoln("1 ", hex.EncodeToString(hash))
	header := pb.Header{
		Method: rootdirectory.MethodUploadBlob,
		//Hash:    hash,
		RawSize: uint64(size),
	}
	rawHeader, err := proto.Marshal(&header)
	if err != nil {
		fmt.Printf("Marshal err:%v\n", err)
		return "", err
	}
	headerLen := uint64(len(rawHeader))
	err = binary.Write(conn, binary.LittleEndian, headerLen)
	if err != nil {
		fmt.Printf("Write headerLen err:%v\n", err)
		return "", err
	}
	_, err = conn.Write(rawHeader)
	if err != nil {
		fmt.Printf("Write rawHeader err:%v\n", err)
		return "", err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	//_, err = f.Seek(0, io.SeekStart)
	//if err != nil {
	//	return "", err
	//}
	defer f.Close()
	_, err = io.Copy(conn, f)
	if err != nil {
		fmt.Printf("Write blob err:%v\n", err)
		return "", err
	}
	var hashLen uint8
	err = binary.Read(conn, binary.LittleEndian, &hashLen)
	if err != nil {
		fmt.Printf("Read hashLen err:%v\n", err)
		return "", err
	}
	hash2 := make([]byte, hashLen)
	_, err = conn.Read(hash2)
	if err != nil {
		fmt.Printf("Read hash err:%v\n", err)
		return "", err
	}
	conn.Close()
	logrus.Infoln("2 ", string(hash2))
	return string(hash2), nil
}

func (c *BackUpCtx) uploadDir(infos []*pb.FileInfo) (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:9877")
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return "", err
	}
	bytes, err := proto.Marshal(&pb.FileInfos{Info: infos})
	if err != nil {
		fmt.Printf("Marshal err:%v\n", err)
		return "", err
	}
	header := pb.Header{
		Method:  rootdirectory.MethodUploadTree,
		RawSize: uint64(len(bytes)),
	}
	rawHeader, err := proto.Marshal(&header)
	if err != nil {
		fmt.Printf("Marshal err:%v\n", err)
		return "", err
	}
	headerLen := uint64(len(rawHeader))
	err = binary.Write(conn, binary.LittleEndian, headerLen)
	if err != nil {
		fmt.Printf("Write headerLen err:%v\n", err)
		return "", err
	}
	_, err = conn.Write(rawHeader)
	if err != nil {
		fmt.Printf("Write rawHeader err:%v\n", err)
		return "", err
	}
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Printf("Write blob err:%v\n", err)
		return "", err
	}
	var hashLen uint8
	err = binary.Read(conn, binary.LittleEndian, &hashLen)
	if err != nil {
		fmt.Printf("Read hashLen err:%v\n", err)
		return "", err
	}
	hash := make([]byte, hashLen)
	_, err = conn.Read(hash)
	if err != nil {
		fmt.Printf("Read hash err:%v\n", err)
		return "", err
	}
	conn.Close()
	logrus.Infoln("2 ", string(hash))
	return string(hash), nil
}

func (c *BackUpCtx) upload(fn func(context.Context, pb.KoalaFSClient) (string, error)) (string, error) {
	conn, err := grpc.Dial(c.host, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Errorf("Dial")
		return "", err
	}
	defer conn.Close()
	client := pb.NewKoalaFSClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "kfs-mount", "backup")
	return fn(ctx, client)
}
