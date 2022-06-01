package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
)

func process(kfsCore *core.KFS, conn net.Conn) {
	reader := bufio.NewReader(conn)

	hashBytes := make([]byte, 256/8)
	err := binary.Read(reader, binary.LittleEndian, hashBytes)
	if err != nil {
		println("hashBytes", err)
		return
	}
	hash := hex.EncodeToString(hashBytes)
	println("hash", hash)

	var size int64
	err = binary.Read(reader, binary.LittleEndian, &size)
	if err != nil {
		println("size", err)
		return
	}
	println("size", size)

	exist, err := kfsCore.S.WriteFn(hash, func(f io.Writer, hasher io.Writer) error {
		_, err = conn.Write([]byte{0}) // not exist
		if err != nil {
			return err
		}
		w := io.MultiWriter(f, hasher)
		_, err = io.CopyN(w, conn, size)
		return err
	})
	if err != nil {
		println("WriteFn", err)
		return
	}
	if exist {
		_, err = conn.Write([]byte{1})
		if err != nil {
			println("exist", err)
		}
		return
	}

	f := sqlite.NewFile(hash, uint64(size))
	err = kfsCore.Db.WriteFile(context.Background(), f)
	if err != nil {
		println("Db.WriteFile", err)
		return
	}

	_, err = conn.Write([]byte{0})
	if err != nil {
		println("code", err)
		return
	}
}

func SocketServer(kfsRoot string, portString string) error {
	kfsCore, _, err := core.New(kfsRoot)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", "127.0.0.1:"+portString)
	if err != nil {
		return err
	}
	println("GRPC listening on", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			println("accept failed", err)
			continue
		}
		go process(kfsCore, conn)
	}
}
