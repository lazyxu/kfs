package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/lazyxu/kfs/db/gosqlite"

	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/server"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	kfsRootStr      = "kfs-root"
	SocketServerStr = "socket-server"
	WebServerStr    = "web-server"
)

func init() {
	rootCmd.PersistentFlags().String(SocketServerStr, "1124", "socket server port")
	rootCmd.PersistentFlags().String(WebServerStr, "1123", "web server port")
}

//go:embed build/*
var build embed.FS

var rootCmd = &cobra.Command{
	Use:   "kfs-server",
	Short: "Kfs is file system used to backup files.",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}()
		kfsRoot := args[0]
		socketPortString := cmd.Flag(SocketServerStr).Value.String()
		socketPort, err := strconv.Atoi(socketPortString)
		if err != nil {
			return
		}
		if socketPort != 0 && socketPort < 1024 || socketPort > 65535 {
			err = errors.New("socketPort should be between 1024 and 65535, actual " + socketPortString)
			return
		}
		webPortString := cmd.Flag(WebServerStr).Value.String()
		webPort, err := strconv.Atoi(webPortString)
		if err != nil {
			return
		}
		if webPort != 0 && webPort < 1024 || webPort > 65535 {
			err = errors.New("webPort should be between 1024 and 65535, actual " + webPortString)
			return
		}
		//kfsCore, err := core.New(mysql.FuncNew("root:12345678@/kfs?parseTime=true&multiStatements=true"), storage.FuncNew(kfsRoot, storage.NewStorage1))
		kfsCore, err := core.New(gosqlite.FuncNew("kfs.db"), storage.FuncNew(kfsRoot, storage.NewStorage1))
		if err != nil {
			return
		}
		defer kfsCore.Close()
		_, err = kfsCore.Checkout(context.Background(), "master")
		if err != nil {
			return
		}
		viper.Set(kfsRootStr, kfsRoot)
		err = viper.WriteConfig()
		if err != nil {
			return
		}
		go func() {
			lis, err := net.Listen("tcp", "0.0.0.0:"+socketPortString)
			if err != nil {
				panic(err)
			}
			println("Socket server listening at:", lis.Addr().String())
			err = server.Socket(lis, kfsCore)
			if err != nil {
				panic(err)
			}
		}()
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsHandler(w, r, kfsCore)
		})
		http.Handle("/", http.FileServer(AddPrefix(http.FS(build), "build")))
		lis, err := net.Listen("tcp", "0.0.0.0:"+webPortString)
		if err != nil {
			panic(err)
		}
		println("Web server listening at:", lis.Addr().String())
		err = http.Serve(lis, nil)
		if err != nil {
			panic(err)
		}
	},
}

type Dir struct {
	fs     http.FileSystem
	prefix string
}

func AddPrefix(fs http.FileSystem, prefix string) http.FileSystem {
	return Dir{fs, prefix}
}

func (d Dir) Open(name string) (http.File, error) {
	return d.fs.Open(path.Clean(d.prefix + name))
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
