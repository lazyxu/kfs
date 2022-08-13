package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/rpcutil"
)

func handleCat(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
	branchName, err := rpcutil.ReadString(conn)
	if err != nil {
		return err
	}
	filePath, err := rpcutil.ReadString(conn)
	if err != nil {
		return err
	}
	println(branchName, filePath)
	ctx := context.Background()
	hash, mode, err := kfsCore.Db.GetFileHashMode(ctx, branchName, core.FormatPath(filePath))
	if err != nil {
		return err
	}
	if !mode.IsRegular() {
		err = fmt.Errorf("invalid mode: %x", mode)
		return err
	}
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = rpcutil.WriteOK(conn)
	if err != nil {
		println(conn.RemoteAddr().String(), "code", err.Error())
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, rc.Size())
	if err != nil {
		return err
	}
	_, err = io.CopyN(conn, rc, rc.Size())
	if err != nil {
		return err
	}
	return nil
}
