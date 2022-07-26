package client

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func (fs *RpcFs) Cat(ctx context.Context, branchName string, srcPath string, fn func(r io.Reader, size int64) error) (err error) {
	conn, err := net.Dial("tcp", fs.SocketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	err = rpcutil.WriteCommandType(conn, rpcutil.CommandCat)
	if err != nil {
		return
	}
	err = rpcutil.WriteString(conn, branchName)
	if err != nil {
		return
	}
	err = rpcutil.WriteString(conn, srcPath)
	if err != nil {
		return
	}

	code, errMsg, err := rpcutil.ReadExit(conn)
	if err != nil {
		return
	}
	if code != rpcutil.EOK {
		err = errors.New(errMsg)
		return
	}
	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		return
	}
	err = fn(conn, size)
	if err != nil {
		return
	}
	return
}
