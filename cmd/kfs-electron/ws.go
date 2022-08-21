package main

import (
	"fmt"
	"io"
	"net"

	"github.com/lazyxu/kfs/rpc/server"

	"github.com/gorilla/websocket"
)

func ToAddrReadWriteCloser(c *websocket.Conn) server.AddrReadWriteCloser {
	return &rwc{c: c}
}

type rwc struct {
	c *websocket.Conn
}

func (c *rwc) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *rwc) Write(p []byte) (int, error) {
	err := c.c.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *rwc) Read(p []byte) (int, error) {
	var r io.Reader
	if r == nil {
		var err error
		_, r, err = c.c.NextReader()
		if err != nil {
			return 0, err
		}
	}
	n, err := r.Read(p)
	if err == io.EOF {
		err = nil
	}
	if n != len(p) {
		panic(fmt.Errorf("invalid read count, expected %d, actual %d", len(p), n))
	}
	return n, err
}

func (c *rwc) Close() error {
	return c.c.Close()
}
