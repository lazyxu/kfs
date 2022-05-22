package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"

	"google.golang.org/grpc"

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

var rootCmd = &cobra.Command{
	Use:   "kfs",
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
		_, err = kfsCore.BranchNew(context.Background(), "master")
		if err != nil {
			return
		}
		viper.Set(kfsRootStr, kfsRoot)
		err = viper.WriteConfig()
		if err != nil {
			return
		}
		server := NewKoalaFSServer(kfsRoot)
		lis, err := net.Listen("tcp", "0.0.0.0:"+portString)
		if err != nil {
			return
		}
		s := grpc.NewServer()
		pb.RegisterKoalaFSServer(s, server)
		println("Listening on", lis.Addr().String())
		err = s.Serve(lis)
		if err != nil {
			return
		}
	},
}
