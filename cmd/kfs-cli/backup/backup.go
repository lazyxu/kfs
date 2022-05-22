package backup

import (
	"fmt"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "backup",
	Example: "kfs-cli backup . -p /test",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runBackup,
}

func init() {
	Cmd.PersistentFlags().String(PathStr, "p", "branch description")
}

func runBackup(cmd *cobra.Command, args []string) {
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

	backupPath := args[0]
	dstPath := cmd.Flag(PathStr).Value.String()

	switch serverType {
	case ServerTypeLocal:
		err = local(cmd.Context(), serverAddr, branchName, dstPath, backupPath)
	case ServerTypeRemote:
		err = remote(cmd.Context(), serverAddr, branchName, dstPath, backupPath)
	default:
		err = InvalidServerType
	}
}
