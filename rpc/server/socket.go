package server

import (
	"io"
	"net"

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	"github.com/lazyxu/kfs/core"
)

type AddrReadWriteCloser interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}

type CommandHandler func(kfsCore *core.KFS, conn AddrReadWriteCloser) error

var commandHandlers = make(map[rpcutil.CommandType]CommandHandler)

func registerCommand(commandType rpcutil.CommandType, handler CommandHandler) {
	commandHandlers[commandType] = func(kfsCore *core.KFS, conn AddrReadWriteCloser) error {
		err := handler(kfsCore, conn)
		if e, ok := err.(*rpcutil.UnexpectedError); ok {
			return e
		}
		if err != nil {
			return rpcutil.WriteInvalid(conn, err)
		}
		return rpcutil.WriteOK(conn)
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
		if e, ok := err.(*rpcutil.UnexpectedError); ok && e.Err == io.EOF {
			conn.Close()
			return
		}
		if err != nil {
			println("commandType", commandType, err.Error())
			conn.Close()
			return
		}
		if handler, ok := commandHandlers[commandType]; ok {
			e := handler(kfsCore, conn)
			if e != nil {
				println(e.Error())
				conn.Close()
				return
			}
		} else {
			println("invalid commandType", commandType)
		}
	}
}

func init() {
	registerCommand(rpcutil.CommandPing, func(kfsCore *core.KFS, conn AddrReadWriteCloser) (err error) {
		_, err = conn.Write([]byte{0})
		return rpcutil.UnexpectedIfError(err)
	})
	registerCommand(rpcutil.CommandReset, handleReset)
	registerCommand(rpcutil.CommandList, handleList)
	registerCommand(rpcutil.CommandUpload, handleUpload)
	registerCommand(rpcutil.CommandTouch, handleTouch)
	registerCommand(rpcutil.CommandDownload, handleDownload)
	registerCommand(rpcutil.CommandCat, handleCat)
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
