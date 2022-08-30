package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/lazyxu/kfs/db/gosqlite"
	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/stretchr/testify/assert"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/rpc/server"
)

const kfsRoot = "test-root-dir"

var (
	socketPort int
)

func initServer() error {
	//kfsCore, err := core.New(mysql.FuncNew("root:12345678@/kfs?parseTime=true&multiStatements=true"), storage.FuncNew(kfsRoot, storage.NewStorage1))
	kfsCore, err := core.New(gosqlite.FuncNew("kfs.db"), storage.FuncNew(kfsRoot, storage.NewStorage1))
	if err != nil {
		return err
	}

	socketLis, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	socketPort = socketLis.Addr().(*net.TCPAddr).Port
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
	cmd := rootCmd()
	cmd.SetArgs(args)
	outBuffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	cmd.SetOut(outBuffer)
	cmd.SetErr(errBuffer)
	err := cmd.Execute()
	if err != nil {
		return "", "", err
	}
	return outBuffer.String(), errBuffer.String(), nil
}

func exec(t *testing.T, args []string) (string, string) {
	stdout, stderr, err := execute(args)
	assert.Nil(t, err)
	return stdout, stderr
}
