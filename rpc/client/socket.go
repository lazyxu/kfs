package client

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"os"

	"github.com/pierrec/lz4"
)

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

	//println(conn.LocalAddr().String(), 1)
	length := len(h.encoder)
	header := make([]byte, length+2)
	header[0] = 1
	copy(header[1:], h.encoder)
	header[length+1] = 0
	_, err = conn.Write(header)
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
	//println(conn.LocalAddr().String(), filePath, "hash", hash)

	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}
	//println(conn.LocalAddr().String(), filePath, "size", size)

	var exist bool
	err = binary.Read(conn, binary.LittleEndian, &exist)
	if err != nil {
		return err
	}
	//println(conn.LocalAddr().String(), filePath, "exist", exist)
	if exist {
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if h.encoder == "lz4" {
		w := lz4.NewWriter(conn)
		_, err = io.CopyN(w, f, int64(size))
		if err != nil {
			w.Flush()
			return err
		}
		w.Flush()
	} else {
		_, err = io.CopyN(conn, f, int64(size))
		if err != nil {
			return err
		}
	}

	//println(conn.LocalAddr().String(), filePath, "CopyN")
	var code int8
	err = binary.Read(conn, binary.LittleEndian, &code)
	if err != nil {
		return err
	}
	//println(conn.LocalAddr().String(), filePath, "code", code)
	return nil
}
