package client

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"

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

func (h *uploadHandlers) uploadFile(filePath string, hash string, size uint64) (err error) {
	c, err := h.p.Get()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			h.p.Close(c)
			return
		}
		err = h.p.Put(c)
	}()
	conn := c.(net.Conn)

	p := &Process{
		conn:      conn,
		filePath:  filePath,
		size:      size,
		stackSize: -1,
	}
	p.label = "hash"
	h.ch <- p

	_, err = conn.Write([]byte{1})
	if err != nil {
		return err
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	_, err = conn.Write(hashBytes)
	if err != nil {
		return err
	}

	p.label = "size"
	h.ch <- p
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	p.label = "exist?"
	h.ch <- p
	var exist bool
	err = binary.Read(conn, binary.LittleEndian, &exist)
	if err != nil {
		return err
	}
	p.label = fmt.Sprintf("exist=%t", exist)
	h.ch <- p

	if exist {
		return nil
	}

	p.label = fmt.Sprintf("encoder=%s", h.encoder)
	h.ch <- p
	length := len(h.encoder)
	header := make([]byte, length+1)
	copy(header[:], h.encoder)
	header[length] = 0
	_, err = conn.Write(header)
	if err != nil {
		return err
	}

	p.label = "copyFile"
	h.ch <- p
	err = h.copyFile(conn, filePath, int64(size))
	if err != nil {
		return err
	}

	p.label = "code?"
	h.ch <- p
	var code int8
	err = binary.Read(conn, binary.LittleEndian, &code)
	if err != nil {
		return err
	}
	p.label = fmt.Sprintf("code=%d", code)
	h.ch <- p

	return nil
}
