package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/pierrec/lz4"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
)

type AddrReadWriteCloser interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}

type CommandHandler func(kfsCore *core.KFS, conn AddrReadWriteCloser)

var commandHandlers = make(map[rpcutil.CommandType]CommandHandler)

func registerCommand(commandType rpcutil.CommandType, handler func(kfsCore *core.KFS, conn AddrReadWriteCloser) error) {
	commandHandlers[commandType] = func(kfsCore *core.KFS, conn AddrReadWriteCloser) {
		err := handler(kfsCore, conn)
		if err != nil {
			rpcutil.WriteInvalid(conn, err)
			return
		}
		rpcutil.WriteOK(conn)
	}
}

func Process(kfsCore *core.KFS, conn AddrReadWriteCloser) {
	println(conn.RemoteAddr().String(), "Process")

	for {
		commandType, err := rpcutil.ReadCommandType(conn)
		if err == io.EOF || websocket.IsUnexpectedCloseError(err) {
			conn.Close()
			return
		}
		if err != nil {
			println("commandType", commandType, err.Error())
			conn.Close()
			return
		}
		if handler, ok := commandHandlers[commandType]; ok {
			handler(kfsCore, conn)
		} else {
			println("invalid commandType", commandType)
		}
	}
}

func init() {
	registerCommand(rpcutil.CommandPing, func(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
		conn.Write([]byte{0})
		return nil
	})
	registerCommand(rpcutil.CommandReset, func(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
		branchName, err := rpcutil.ReadString(conn)
		if err != nil {
			return err
		}
		err = kfsCore.Reset(context.TODO(), branchName)
		if err != nil {
			println(conn.RemoteAddr().String(), "Reset", err.Error())
			return err
		}
		return nil
	})
	registerCommand(rpcutil.CommandUpload, handleUpload)
	registerCommand(rpcutil.CommandTouch, handleTouch)
	registerCommand(rpcutil.CommandList, handleList)
	registerCommand(rpcutil.CommandDownload, handleDownload)
	registerCommand(rpcutil.CommandCat, handleCat)
}

func handleUpload(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

	hashBytes := make([]byte, 256/8)
	err := binary.Read(conn, binary.LittleEndian, hashBytes)
	if err != nil {
		println(conn.RemoteAddr().String(), "hashBytes", err.Error())
		return err
	}
	hash := hex.EncodeToString(hashBytes)
	println("hash", hash)

	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		println(conn.RemoteAddr().String(), "size", err.Error())
		return err
	}
	println(conn.RemoteAddr().String(), "size", size)

	exist, err := kfsCore.S.WriteFn(hash, func(f io.Writer, hasher io.Writer) error {
		_, err = conn.Write([]byte{1}) // not exist
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
		return err
	}
	if exist {
		return nil
	}

	f := sqlite.NewFile(hash, uint64(size))
	err = kfsCore.Db.WriteFile(context.Background(), f)
	if err != nil {
		println(conn.RemoteAddr().String(), "Db.WriteFile", err.Error())
		return err
	}
	return nil
}

func Socket(listener net.Listener, kfsCore *core.KFS) error {
	println("Socket listening on", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("accept failed", err)
			continue
		}
		go Process(kfsCore, conn)
	}
}
