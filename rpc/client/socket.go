package client

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/pierrec/lz4"
)

func (h *uploadHandlers) copyFile(conn net.Conn, filePath string, size int64) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if h.encoder == "lz4" {
		w := lz4.NewWriter(conn)
		_, err = io.CopyN(w, f, size)
		if err != nil {
			return err
		}
		defer w.Flush()
	} else {
		_, err = io.CopyN(conn, f, size)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *uploadHandlers) getSizeAndCalHash(filePath string, p *Process) (sqlite.File, error) {
	if h.verbose {
		p.label = "stat?"
		h.ch <- p
	}
	f, err := os.Open(filePath)
	if err != nil {
		return sqlite.File{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return sqlite.File{}, err
	}
	if h.verbose {
		p.label = "hash?"
		p.size = uint64(info.Size())
		h.ch <- p
	}
	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return sqlite.File{}, err
	}
	return sqlite.NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(info.Size())), nil
}

func (h *uploadHandlers) uploadFile(ctx context.Context, index int, filePath string) (file sqlite.File, err error) {
	var p *Process
	if h.verbose {
		p = &Process{
			index:     index,
			filePath:  filePath,
			stackSize: -1,
		}
		p.label = "start"
		h.ch <- p
	}
	file, err = h.getSizeAndCalHash(filePath, p)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			h.conns[index].Close()
			h.BeforeFileHandler(ctx, index)
			return
		}
	}()
	conn := h.conns[index]

	err = rpcutil.WriteCommandType(conn, rpcutil.CommandUpload)
	if err != nil {
		return
	}

	hashBytes, err := hex.DecodeString(file.Hash())
	if err != nil {
		return
	}
	_, err = conn.Write(hashBytes)
	if err != nil {
		return
	}

	if h.verbose {
		p.label = "size"
		h.ch <- p
	}
	err = binary.Write(conn, binary.LittleEndian, file.Size())
	if err != nil {
		return
	}

	if h.verbose {
		p.label = "exist?"
		h.ch <- p
	}
	var exist bool
	err = binary.Read(conn, binary.LittleEndian, &exist)
	if err != nil {
		return
	}

	if exist {
		if h.verbose {
			p.label = fmt.Sprintf("exist")
			h.ch <- p
		}
		return
	}

	if h.verbose {
		p.label = fmt.Sprintf("e=%s", h.encoder)
		h.ch <- p
	}
	err = rpcutil.WriteString(conn, h.encoder)
	if err != nil {
		return
	}

	if h.verbose {
		p.label = "copyFile"
		h.ch <- p
	}
	err = h.copyFile(conn, filePath, int64(file.Size()))
	if err != nil {
		return
	}

	if h.verbose {
		p.label = "code?"
		h.ch <- p
	}
	code, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if h.verbose {
		p.label = fmt.Sprintf("code=%d", code)
		h.ch <- p
		if code != rpcutil.EOK {
			p.err = errors.New(errMsg)
			h.ch <- p
		}
	}

	return
}
