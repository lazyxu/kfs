package main

import (
	"fmt"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:     "reset",
	Example: "kfs-cli reset",
	Args:    cobra.RangeArgs(0, 0),
	Run:     runReset,
}

func runReset(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(err)
	}()
	serverType := viper.GetString(ServerTypeStr)
	serverAddr := viper.GetString(ServerAddrStr)
	branchName := viper.GetString(BranchNameStr)
	fmt.Printf("%s: %s\n", ServerTypeStr, serverType)
	fmt.Printf("%s: %s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s: %s\n", BranchNameStr, branchName)

	err = withFS(serverType, serverAddr, func(fs core.FS) error {
		return fs.Reset(cmd.Context())
	})
}
