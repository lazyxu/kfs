package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleCat(kfsCore *core.KFS, conn net.Conn) {
	var err error
	defer func() {
		if err != nil {
			rpcutil.WriteInvalid(conn, err)
		}
	}()
	branchName, err := rpcutil.ReadString(conn)
	if err != nil {
		return
	}
	filePath, err := rpcutil.ReadString(conn)
	if err != nil {
		return
	}
	println(branchName, filePath)
	ctx := context.Background()
	hash, mode, err := kfsCore.Db.GetFileHashMode(ctx, branchName, core.FormatPath(filePath))
	if err != nil {
		return
	}
	if !mode.IsRegular() {
		err = fmt.Errorf("invalid mode: %x", mode)
		return
	}
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	err = rpcutil.WriteOK(conn)
	if err != nil {
		println(conn.RemoteAddr().String(), "code", err.Error())
		return
	}
	err = binary.Write(conn, binary.LittleEndian, rc.Size())
	if err != nil {
		return
	}
	_, err = io.CopyN(conn, rc, rc.Size())
	if err != nil {
		return
	}
	err = rpcutil.WriteOK(conn)
}
