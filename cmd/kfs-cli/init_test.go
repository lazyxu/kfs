package main

import (
	"bytes"
	"fmt"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/gosqlite"
	"github.com/lazyxu/kfs/db/mysql"
	storage "github.com/lazyxu/kfs/storage/local"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/rpc/server"
)

const kfsRoot = "test-root-dir"

var (
	socketPort int
)

func initServer() error {
	storageType := os.Getenv("kfs_test_storage_type")
	var newStorage func() (dao.Storage, error)
	if storageType == "0" {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage0)
	} else if storageType == "1" {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage1)
	} else if storageType == "2" {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage2)
	} else if storageType == "3" {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage3)
	} else if storageType == "4" {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage4)
	} else {
		newStorage = dao.StorageNewFunc(kfsRoot, storage.NewStorage5)
	}

	databaseType := os.Getenv("kfs_test_database_type")
	var newDatabase func() (dao.Database, error)
	if databaseType == "mysql" {
		newDatabase = dao.DatabaseNewFunc("root:12345678@/kfs?parseTime=true&multiStatements=true", mysql.New)
	} else {
		newDatabase = dao.DatabaseNewFunc("kfs.db", gosqlite.New)
	}

	kfsCore, err := core.New(newDatabase, newStorage)
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
