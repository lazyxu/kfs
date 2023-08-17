package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/gosqlite"
	"github.com/lazyxu/kfs/db/mysql"
	storage "github.com/lazyxu/kfs/storage/local"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/server"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
)

const (
	SocketServerStr   = "socket-server"
	WebServerStr      = "web-server"
	DatabaseTypeStr   = "database-type"
	DataSourceNameStr = "data-source-name"
	StorageTypeStr    = "storage-type"
	StorageDirStr     = "storage-dir"
)

func init() {
	rootCmd.PersistentFlags().String(WebServerStr, "1123", "web server port")
	rootCmd.PersistentFlags().String(SocketServerStr, "1124", "socket server port")
	rootCmd.PersistentFlags().String(DatabaseTypeStr, "sqlite", "sqlite/mysql")
	rootCmd.PersistentFlags().String(DataSourceNameStr, "kfs.db", "data source name")
	rootCmd.PersistentFlags().String(StorageTypeStr, "1", "storage type, [0, 5]")
	rootCmd.PersistentFlags().String(StorageDirStr, "tmp", "storage dir path")
}

//go:embed build/*
var build embed.FS

func getStorageByType(typ string) (func(string) (dao.Storage, error), error) {
	switch typ {
	case "0":
		return storage.NewStorage0, nil
	case "1":
		return storage.NewStorage1, nil
	case "2":
		return storage.NewStorage2, nil
	case "3":
		return storage.NewStorage3, nil
	case "4":
		return storage.NewStorage4, nil
	case "5":
		return storage.NewStorage5, nil
	}
	return nil, fmt.Errorf("no such storage type: %s", typ)
}

func getDatabaseByType(typ string) (func(string) (dao.Database, error), error) {
	switch typ {
	case "sqlite":
		return gosqlite.New, nil
	case "mysql":
		return mysql.New, nil
	}
	return nil, fmt.Errorf("no such databse type: %s", typ)
}

var kfsCore *core.KFS

var rootCmd = &cobra.Command{
	Use:   "kfs-server",
	Short: "Kfs is file system used to backup files.",
	Args:  cobra.RangeArgs(0, 0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}()

		storageType := cmd.Flag(StorageTypeStr).Value.String()
		storageDir := cmd.Flag(StorageDirStr).Value.String()

		databaseType := cmd.Flag(DatabaseTypeStr).Value.String()
		dataSourceName := cmd.Flag(DataSourceNameStr).Value.String()

		webPortString := cmd.Flag(WebServerStr).Value.String()
		webPort, err := strconv.Atoi(webPortString)
		if err != nil {
			return
		}
		if webPort != 0 && webPort < 1024 || webPort > 65535 {
			err = errors.New("webPort should be between 1024 and 65535, actual " + webPortString)
			return
		}

		socketPortString := cmd.Flag(SocketServerStr).Value.String()
		socketPort, err := strconv.Atoi(socketPortString)
		if err != nil {
			return
		}
		if socketPort != 0 && socketPort < 1024 || socketPort > 65535 {
			err = errors.New("socketPort should be between 1024 and 65535, actual " + socketPortString)
			return
		}

		newStorage, err := getStorageByType(storageType)
		if err != nil {
			return
		}
		newDatabase, err := getDatabaseByType(databaseType)
		if err != nil {
			return
		}

		kfsCore, err = core.New(dao.DatabaseNewFunc(dataSourceName, newDatabase), dao.StorageNewFunc(storageDir, newStorage))
		if err != nil {
			return
		}
		defer kfsCore.Close()
		//GetExifData("375667db5da6ed4017815f864ffe0563182523167ce40448c175298fe6af56d1")
		err = diskUsage()
		if err != nil {
			println("diskUsage:", err.Error())
		}
		//AnalysisFileType(context.TODO())

		go func() {
			// socket
			lis, err := net.Listen("tcp", "0.0.0.0:"+socketPortString)
			if err != nil {
				panic(err)
			}
			println("KFS socket server listening at:", lis.Addr().String())
			err = server.Socket(lis, kfsCore)
			if err != nil {
				panic(err)
			}
		}()
		// web
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsHandler(w, r, kfsCore)
		})
		webServer(webPortString)
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func wsHandler(w http.ResponseWriter, r *http.Request, kfsCore *core.KFS) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	server.Process(kfsCore, ToAddrReadWriteCloser(c))
}
