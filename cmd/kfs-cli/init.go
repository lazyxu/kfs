package main

import (
	"fmt"

	"github.com/lazyxu/kfs/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use: "init",
	Example: `
kfs-cli init -t local -b master ./tmp
kfs-cli init -t remote -b master localhost:1123
`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runInit,
}

func init() {
	initCmd.PersistentFlags().StringP(ServerTypeStr, "t", "remote", "local/remote")
	initCmd.PersistentFlags().StringP(BranchNameStr, "b", "master", "")
}

func runInit(cmd *cobra.Command, args []string) {
	var err error
	serverType := cmd.Flag(ServerTypeStr).Value.String()
	grpcServerAddr := args[0]
	socketServerAddr := "localhost:1124"
	branchName := cmd.Flag(BranchNameStr).Value.String()
	defer func() {
		viper.Set(ServerTypeStr, serverType)
		viper.Set(GrpcServerAddrStr, grpcServerAddr)
		viper.Set(SocketServerAddrStr, socketServerAddr)
		viper.Set(BranchNameStr, branchName)
		err = viper.WriteConfig()
		ExitWithError(err)
		fmt.Printf("%s: %s\n", ServerTypeStr, serverType)
		fmt.Printf("%s: %s\n", GrpcServerAddrStr, grpcServerAddr)
		fmt.Printf("%s: %s\n", SocketServerAddrStr, socketServerAddr)
		fmt.Printf("%s: %s\n", BranchNameStr, branchName)
	}()
	defer func() {
		ExitWithError(err)
	}()

	fs := &client.RpcFs{
		GrpcServerAddr:   grpcServerAddr,
		SocketServerAddr: socketServerAddr,
	}

	_, err = fs.Checkout(cmd.Context(), branchName)
}
