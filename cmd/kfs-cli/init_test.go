package main

import (
	"bytes"
	"fmt"
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
	stdout, stderr, err := execute([]string{"init",
		"--grpc-server", "localhost:" + strconv.Itoa(grpcPort),
		"--socket-server", "localhost:" + strconv.Itoa(socketPort),
		"--config-file", ".kfs.json"})
	if err != nil {
		panic(err)
	}
	if stdout != "" {
		panic(fmt.Errorf("init expected \"\", actual \"%s\"", stdout))
	}
	if stderr != "" {
		panic(fmt.Errorf("init expected \"\", actual \"%s\"", stderr))
	}
}

func execute(args []string) (string, string, error) {
	rootCmd.SetArgs(args)
	outBuffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	rootCmd.SetOut(outBuffer)
	rootCmd.SetErr(errBuffer)
	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}
	return outBuffer.String(), errBuffer.String(), nil
}

func exec(t *testing.T, args []string) (string, string) {
	stdout, stderr, err := execute(args)
	if err != nil {
		t.Error(err)
		return stdout, stderr
	}
	return stdout, stderr
}
