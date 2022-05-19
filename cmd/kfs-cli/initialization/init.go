package initialization

import (
	"fmt"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use: "init",
	Example: `
kfs-cli init -t local -b master ./tmp
kfs-cli init -t remote -b master localhost:1123
`,
	Args: cobra.RangeArgs(1, 1),
	Run:  runInit,
}

func init() {
	Cmd.PersistentFlags().StringP(ServerTypeStr, "t", "remote", "local/remote")
	Cmd.PersistentFlags().StringP(BranchNameStr, "b", "master", "")
	Cmd.PersistentFlags().String(DescriptionStr, "", "branch description")
}

func runInit(cmd *cobra.Command, args []string) {
	serverType := cmd.Flag(ServerTypeStr).Value.String()
	serverAddr := args[0]
	branchName := cmd.Flag(BranchNameStr).Value.String()
	description := cmd.Flag(DescriptionStr).Value.String()
	defer func() {
		viper.Set(ServerTypeStr, serverType)
		viper.Set(ServerAddrStr, serverAddr)
		viper.Set(BranchNameStr, branchName)
		err := viper.WriteConfig()
		ExitWithError(err)
		fmt.Printf("%s=%s\n", ServerTypeStr, serverType)
		fmt.Printf("%s=%s\n", ServerAddrStr, serverAddr)
		fmt.Printf("%s=%s\n", BranchNameStr, branchName)
		fmt.Printf("%s=%s\n", DescriptionStr, description)
	}()
	var err error
	defer ExitWithError(err)
	switch serverType {
	case ServerTypeLocal:
		serverAddr, err = local(serverAddr, branchName, description)
	case ServerTypeRemote:
		err = remote(serverAddr, branchName, description)
	default:
		err = InvalidServerType
	}
}
