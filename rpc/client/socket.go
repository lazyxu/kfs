package client

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"io"
	"net"
	"os"

	"github.com/lazyxu/kfs/rpc/rpcutil"

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

func (h *uploadHandlers) getSizeAndCalHash(filePath string, p *core.Process) (dao.File, error) {
	if h.verbose {
		p.Label = "stat?"
		h.uploadProcess.Show(p)
	}
	f, err := os.Open(filePath)
	if err != nil {
		return dao.File{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return dao.File{}, err
	}
	if h.verbose {
		p.Label = "hash?"
		p.Size = uint64(info.Size())
		h.uploadProcess.Show(p)
	}
	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return dao.File{}, err
	}
	return dao.NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(info.Size())), nil
}

func (h *uploadHandlers) uploadFile(ctx context.Context, index int, filePath string) (file dao.File, err error) {
	var p *core.Process
	if h.verbose {
		p = &core.Process{
			Index:     index,
			FilePath:  filePath,
			StackSize: -1,
		}
		p.Label = "start"
		h.uploadProcess.Show(p)
	}
	file, err = h.getSizeAndCalHash(filePath, p)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			h.conns[index].Close()
			h.StartWorker(ctx, index)
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
		p.Label = "size"
		h.uploadProcess.Show(p)
	}
	err = binary.Write(conn, binary.LittleEndian, file.Size())
	if err != nil {
		return
	}

	if h.verbose {
		p.Label = "exist?"
		h.uploadProcess.Show(p)
	}
	var notExist bool
	err = binary.Read(conn, binary.LittleEndian, &notExist)
	if err != nil {
		return
	}

	if !notExist {
		if h.verbose {
			p.Label = fmt.Sprintf("exist")
			h.uploadProcess.Show(p)
		}
		return
	}

	if h.verbose {
		p.Label = fmt.Sprintf("e=%s", h.encoder)
		h.uploadProcess.Show(p)
	}
	err = rpcutil.WriteString(conn, h.encoder)
	if err != nil {
		return
	}

	if h.verbose {
		p.Label = "copyFile"
		h.uploadProcess.Show(p)
	}
	err = h.copyFile(conn, filePath, int64(file.Size()))
	if err != nil {
		return
	}

	if h.verbose {
		p.Label = "code?"
		h.uploadProcess.Show(p)
	}
	code, errMsg, err := rpcutil.ReadStatus(conn)
	if err != nil {
		return
	}
	if h.verbose {
		p.Label = fmt.Sprintf("code=%d", code)
		h.uploadProcess.Show(p)
		if code != rpcutil.EOK {
			p.Err = errors.New(errMsg)
			h.uploadProcess.Show(p)
		}
	}

	return
}
