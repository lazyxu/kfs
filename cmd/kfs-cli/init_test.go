package main

import (
	"bytes"
	"net"
	"strconv"
	"strings"
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
}

func execute(args []string) (string, error) {
	rootCmd.SetArgs(args)
	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	err := rootCmd.Execute()
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func Test_init_ls(t *testing.T) {
	initOutput, err := execute([]string{"init",
		"--grpc-server", "localhost:" + strconv.Itoa(grpcPort),
		"--socket-server", "localhost:" + strconv.Itoa(socketPort),
		"--config-file", ".kfs.json"})
	if err != nil {
		t.Error(err)
	}
	if initOutput != "" {
		t.Errorf("init expected \"\", actual \"%s\"", initOutput)
	}

	lsOutput, err := execute([]string{"ls"})
	if err != nil {
		t.Error(err)
	}
	lsOutput = strings.Trim(lsOutput, "\n")
	if lsOutput != "total 0" {
		t.Errorf("expected \"total 0\", actual \"%s\"", lsOutput)
	}
}
