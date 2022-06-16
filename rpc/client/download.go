package client

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/core"
)

func (fs GRPCFS) Download(ctx context.Context, branchName string, dstPath string, srcPath string, config core.UploadConfig) (filePath string, err error) {
	conn, err := net.Dial("tcp", "127.0.0.1:1124")
	if err != nil {
		return
	}
	defer conn.Close()

	err = rpcutil.WriteCommandType(conn, rpcutil.CommandDownload)
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

	filePath, err = getFilePath(dstPath)
	if err != nil {
		return
	}
	for {
		var mode os.FileMode
		err = binary.Read(conn, binary.LittleEndian, &mode)
		if err != nil {
			return
		}
		if mode == 0 {
			break
		}
		var relPath string
		relPath, err = rpcutil.ReadString(conn)
		if err != nil {
			return
		}
		curPath := filePath
		if relPath != "" {
			curPath += "/" + relPath
		}
		println("download", relPath, mode.IsDir())
		if mode.IsDir() {
			err = os.Mkdir(curPath, mode)
			if err != nil {
				return
			}
		} else {
			err = writeFile(conn, curPath, mode)
			if err != nil {
				return
			}
		}
	}
	return
}

func getFilePath(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return filePath, nil
	}
	if err != nil {
		return "", err
	}
	for i := 0; ; i++ {
		curPath := filePath + "." + strconv.Itoa(i)
		_, err := os.Stat(curPath)
		if os.IsNotExist(err) {
			return curPath, nil
		}
		if err != nil {
			return "", err
		}
	}
}

func writeFile(r io.Reader, filePath string, mode os.FileMode) error {
	var size int64
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.CopyN(f, r, size)
	if err != nil {
		return err
	}
	err = f.Chmod(mode)
	if err != nil {
		return err
	}
	return nil
}
