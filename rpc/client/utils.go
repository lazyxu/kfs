package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"path/filepath"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"google.golang.org/protobuf/proto"
)

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

func FormatFilename(filename string) string {
	var name = []rune(filepath.Base(filename))
	if len(name) > 10 {
		name = append(name[:10], []rune("..")...)
	}
	return fmt.Sprintf("%-12s", string(name))
}
