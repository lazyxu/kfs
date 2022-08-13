package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/core"
)

func handleDownload(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
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
	err = download(ctx, kfsCore, conn, "", hash, mode)
	if err != nil {
		return err
	}
	mode = 0
	err = binary.Write(conn, binary.LittleEndian, mode)
	if err != nil {
		return err
	}
	return nil
}

func download(ctx context.Context, kfsCore *core.KFS, conn AddrReadWriteCloser, relPath string, hash string, mode os.FileMode) error {
	println("download", relPath, hash[:4], mode.IsDir())
	err := binary.Write(conn, binary.LittleEndian, mode)
	if err != nil {
		return err
	}
	err = rpcutil.WriteString(conn, relPath)
	if err != nil {
		return err
	}
	if !mode.IsDir() && !mode.IsRegular() {
		err = fmt.Errorf("invalid mode: %x", mode)
		return err
	}
	if mode.IsDir() {
		dirItems, err := kfsCore.Db.ListByHash(ctx, hash)
		if err != nil {
			return err
		}
		for _, item := range dirItems {
			err = download(ctx, kfsCore, conn, relPath+"/"+item.Name, item.Hash, os.FileMode(item.Mode))
			if err != nil {
				return err
			}
		}
		return nil
	}
	rc, err := kfsCore.S.ReadWithSize(hash)
	defer rc.Close()
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
