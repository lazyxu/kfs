package main

import (
	"fmt"

	"github.com/lazyxu/kfs/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "init",
		Example: `
kfs-cli init -b master --socket-server localhost:1124
`,
		Args: cobra.RangeArgs(0, 0),
		Run:  runInit,
	}
	cmd.PersistentFlags().StringP(BranchNameStr, "b", "master", "default branch name")
	cmd.PersistentFlags().String(SocketServerStr, "localhost:1124", "socket server address")
	return cmd
}

func runInit(cmd *cobra.Command, args []string) {
	var err error
	socketServer := cmd.Flag(SocketServerStr).Value.String()
	branchName := cmd.Flag(BranchNameStr).Value.String()
	defer func() {
		loadConfigFile(cmd)
		verbose := cmd.Flag(VerboseStr).Value.String() != "false"
		viper.Set(SocketServerStr, socketServer)
		viper.Set(BranchNameStr, branchName)
		err = viper.WriteConfig()
		ExitWithError(cmd, err)
		if verbose {
			fmt.Printf("%s: %s\n", SocketServerStr, socketServer)
			fmt.Printf("%s: %s\n", BranchNameStr, branchName)
		}
	}()
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs := &client.RpcFs{
		SocketServerAddr: socketServer,
	}

	_, err = fs.Checkout(cmd.Context(), branchName)
}
