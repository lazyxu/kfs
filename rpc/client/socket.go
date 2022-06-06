package client

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"os"
)

func (v *uploadVisitor) uploadFile(filePath string, hash string, size uint64) (err error) {
	c, err := v.p.Get()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			v.p.Close(c)
			return
		}
		err = v.p.Put(c)
	}()
	conn := c.(net.Conn)

	//println(conn.LocalAddr().String(), 1)
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

	_, err = io.CopyN(conn, f, int64(size))
	if err != nil {
		return err
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
