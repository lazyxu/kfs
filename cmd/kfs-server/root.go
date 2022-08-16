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

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/server"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	kfsRootStr    = "kfs-root"
	backupPathStr = "backup-path"
	branchNameStr = "branch-name"
	pathStr       = "path"
	portStr       = "port"
)

func init() {
	rootCmd.PersistentFlags().StringP(portStr, "p", "0", "grpc port")
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
		portString := cmd.Flag(portStr).Value.String()
		port, err := strconv.Atoi(portString)
		if err != nil {
			return
		}
		if port != 0 && port < 1024 || port > 65535 {
			err = errors.New("port should be between 1024 and 15535, actual " + portString)
			return
		}
		kfsCore, _, err := core.New(kfsRoot)
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
			lis, err := net.Listen("tcp", "0.0.0.0:1124")
			if err != nil {
				panic(err)
			}
			err = server.Socket(lis, kfsCore)
			if err != nil {
				panic(err)
			}
		}()
		go func() {
			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				wsHandler(w, r, kfsCore)
			})
			http.Handle("/", http.FileServer(AddPrefix(http.FS(build), "build")))
			log.Fatal(http.ListenAndServe("0.0.0.0:1125", nil))
		}()
		lis, err := net.Listen("tcp", "0.0.0.0:1123")
		if err != nil {
			panic(err)
		}
		err = server.Grpc(lis, kfsCore)
		if err != nil {
			return
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
