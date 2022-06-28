package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/pierrec/lz4"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
)

func process(kfsCore *core.KFS, conn net.Conn) {
	println(conn.RemoteAddr().String(), "process")

	for {
		commandType, err := rpcutil.ReadCommandType(conn)
		if err != nil && err != io.EOF {
			println("commandType", err.Error())
			return
		}
		switch commandType {
		case rpcutil.CommandPing:
			pong(conn)
		case rpcutil.CommandUpload:
			handleUpload(kfsCore, conn)
		case rpcutil.CommandDownload:
			handleDownload(kfsCore, conn)
		default:
			println("no such commandType", commandType)
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

func handleUpload(kfsCore *core.KFS, conn net.Conn) {
	var err error
	defer func() {
		if err != nil {
			rpcutil.WriteErrorExit(conn, err)
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

		encoder, err := rpcutil.ReadString(conn)
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

	err = rpcutil.WriteSuccessExit(conn)
	if err != nil {
		println(conn.RemoteAddr().String(), "code", err.Error())
		return
	}
	println(conn.RemoteAddr().String(), "code", 0)
}

func Socket(listener net.Listener, kfsCore *core.KFS) error {
	println("Socket listening on", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("accept failed", err)
			continue
		}
		go process(kfsCore, conn)
	}
}
