package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
)

func process(kfsCore *core.KFS, conn net.Conn) {
	println(conn.RemoteAddr().String(), "process")

	for {
		var command uint8
		err := binary.Read(conn, binary.LittleEndian, &command)
		if err != nil {
			println("command", err.Error())
			return
		}
		switch command {
		case 0:
			pong(conn)
		case 1:
			handleUploadFile(kfsCore, conn)
		default:
			panic(fmt.Errorf("no such command %d", command))
		}
	}
}

func pong(conn net.Conn) {
	_, err := conn.Write([]byte{0})
	if err != nil {
		println(conn.RemoteAddr().String(), "pong", err)
		return
	}
}

func handleUploadFile(kfsCore *core.KFS, conn net.Conn) {
	var err error
	defer func() {
		if err != nil {
			_, err = conn.Write([]byte{1})
		}
	}()
	reader := bufio.NewReader(conn)

	hashBytes := make([]byte, 256/8)
	err = binary.Read(reader, binary.LittleEndian, hashBytes)
	if err != nil {
		println(conn.RemoteAddr().String(), "hashBytes", err.Error())
		return
	}
	hash := hex.EncodeToString(hashBytes)
	println("hash", hash)

	var size int64
	err = binary.Read(reader, binary.LittleEndian, &size)
	if err != nil {
		println(conn.RemoteAddr().String(), "size", err.Error())
		return
	}
	println(conn.RemoteAddr().String(), "size", size)

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
		println(conn.RemoteAddr().String(), "WriteFn", err.Error())
		return
	}
	if exist {
		_, err = conn.Write([]byte{1})
		if err != nil {
			println(conn.RemoteAddr().String(), "exist", err.Error())
		}
		println(conn.RemoteAddr().String(), "exist")
		return
	}

	f := sqlite.NewFile(hash, uint64(size))
	err = kfsCore.Db.WriteFile(context.Background(), f)
	if err != nil {
		println(conn.RemoteAddr().String(), "Db.WriteFile", err.Error())
		return
	}

	_, err = conn.Write([]byte{0})
	if err != nil {
		println(conn.RemoteAddr().String(), "code", err.Error())
		return
	}
	println(conn.RemoteAddr().String(), "code", 0)
}

func SocketServer(kfsCore *core.KFS, portString string) error {
	lis, err := net.Listen("tcp", "0.0.0.0:"+portString)
	if err != nil {
		return err
	}
	println("Socket listening on", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			println("accept failed", err)
			continue
		}
		go process(kfsCore, conn)
	}
}
