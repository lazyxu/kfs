package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"strings"

	"github.com/pierrec/lz4"

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
			println("no such command", command)
			//panic(fmt.Errorf("no such command %d", command))
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

	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

	hashBytes := make([]byte, 256/8)
	err = binary.Read(conn, binary.LittleEndian, hashBytes)
	if err != nil {
		println(conn.RemoteAddr().String(), "hashBytes", err.Error())
		return
	}
	hash := hex.EncodeToString(hashBytes)
	println("hash", hash)

	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		println(conn.RemoteAddr().String(), "size", err.Error())
		return
	}
	println(conn.RemoteAddr().String(), "size", size)

	exist, err := kfsCore.S.WriteFn(hash, func(f io.Writer, hasher io.Writer) error {
		_, err := conn.Write([]byte{0}) // not exist
		if err != nil {
			return err
		}

		encoder, err := readString(conn)
		println(conn.RemoteAddr().String(), "encoder", len(encoder), encoder)

		w := io.MultiWriter(f, hasher)
		if encoder == "lz4" {
			r := lz4.NewReader(conn)
			_, err = io.CopyN(w, r, size)
		} else {
			_, err = io.CopyN(w, conn, size)
		}
		println(conn.RemoteAddr().String(), "Copy")
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

func readString(r io.Reader) (string, error) {
	var b byte
	var sb strings.Builder
	for {
		err := binary.Read(r, binary.LittleEndian, &b)
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		err = sb.WriteByte(b)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
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
