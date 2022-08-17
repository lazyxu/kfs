package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/lazyxu/kfs/pb"
)

const fileChunkSize = 1024 * 1024

func ReqString(socketServerAddr string, commandType rpcutil.CommandType, req string) (err error) {
	conn, err := net.Dial("tcp", socketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	// write
	err = rpcutil.WriteCommandType(conn, commandType)
	if err != nil {
		return
	}
	err = rpcutil.WriteString(conn, req)
	if err != nil {
		return
	}

	// read
	status, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}
	return nil
}

func ReqResp(socketServerAddr string, commandType rpcutil.CommandType, req proto.Message, resp proto.Message) (err error) {
	conn, err := net.Dial("tcp", socketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	// write
	err = rpcutil.WriteCommandType(conn, commandType)
	if err != nil {
		return
	}
	err = rpcutil.WriteProto(conn, req)
	if err != nil {
		return
	}

	// read
	status, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}

	err = rpcutil.ReadProto(conn, resp)
	if err != nil {
		return
	}

	// exit
	status, errMsg, err = rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}
	return nil
}

func ReqStringResp(socketServerAddr string, commandType rpcutil.CommandType, req string, resp proto.Message) (err error) {
	conn, err := net.Dial("tcp", socketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	// write
	err = rpcutil.WriteCommandType(conn, commandType)
	if err != nil {
		return
	}
	err = rpcutil.WriteString(conn, req)
	if err != nil {
		return
	}

	// read
	status, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}

	err = rpcutil.ReadProto(conn, resp)
	if err != nil {
		return
	}

	// exit
	status, errMsg, err = rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}
	return nil
}

func ReqRespN(socketServerAddr string, commandType rpcutil.CommandType, req proto.Message, resp proto.Message,
	onLength func(int64) error, onDirItem func() error) (err error) {
	conn, err := net.Dial("tcp", socketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	// write
	err = rpcutil.WriteCommandType(conn, commandType)
	if err != nil {
		return
	}
	err = rpcutil.WriteProto(conn, req)
	if err != nil {
		return
	}

	// read
	status, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}

	var n int64
	err = binary.Read(conn, binary.LittleEndian, &n)
	if err != nil {
		return
	}
	err = onLength(n)
	if err != nil {
		return
	}

	for i := int64(0); i < n; i++ {
		err = rpcutil.ReadProto(conn, resp)
		if err != nil {
			return
		}
		err = onDirItem()
		if err != nil {
			return
		}
	}

	// exit
	status, errMsg, err = rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if status != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}
	return nil
}

func getGRPCClient(fs *RpcFs) (*grpc.ClientConn, pb.KoalaFSClient, error) {
	conn, err := grpc.Dial(fs.GrpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewKoalaFSClient(conn)
	return conn, c, nil
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
