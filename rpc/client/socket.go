package client

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"os"
)

func (v *uploadVisitor) getConn() {
}

func (v *uploadVisitor) uploadFile(filePath string, hash string, size uint64) error {
	//conn, err := net.Dial("tcp", "127.0.0.1:1124")
	//if err != nil {
	//	println(err.Error())
	//	return nil
	//}
	//defer conn.Close()
	c := <-v.connCh
	defer func() {
		println("conn 2", filePath, hash, c)
		v.connCh <- c
	}()
	println("conn 1", filePath, hash, c)

	conn, err := net.Dial("tcp", "127.0.0.1:1124")
	if err != nil {
		return err
	}
	defer conn.Close()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	_, err = conn.Write(hashBytes)
	if err != nil {
		return err
	}

	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	var exist bool
	err = binary.Read(conn, binary.LittleEndian, &exist)
	if err != nil {
		return err
	}
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

	var code int8
	err = binary.Read(conn, binary.LittleEndian, &code)
	if err != nil {
		return err
	}
	return nil
}
