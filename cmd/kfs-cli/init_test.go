package main

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/rpc/server"
)

const kfsRoot = "test-root-dir"

var (
	socketPort int
	grpcPort   int
)

func initServer() error {
	kfsCore, _, err := core.New(kfsRoot)
	if err != nil {
		return err
	}
	_, err = kfsCore.Checkout(context.Background(), "master")
	if err != nil {
		return err
	}

	socketLis, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	socketPort = socketLis.Addr().(*net.TCPAddr).Port
	grpcLis, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	grpcPort = grpcLis.Addr().(*net.TCPAddr).Port
	go func() {
		err = server.Grpc(grpcLis, kfsCore)
		if err != nil {
			return
		}
	}()
	go func() {
		err = server.Socket(socketLis, kfsCore)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func init() {
	err := initServer()
	if err != nil {
		panic(err)
	}
}

func TestInit(t *testing.T) {
	rootCmd.SetArgs([]string{"init", "localhost:" + strconv.Itoa(grpcPort), "--config-file", ".kfs.json"})
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	err := rootCmd.Execute()
	if err != nil {
		t.Error(err)
	}
}
