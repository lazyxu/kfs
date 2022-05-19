package list

import (
	"fmt"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func printHeader(total int) {
	fmt.Printf("total %d\n", total)
	if total != 0 {
		fmt.Printf("mode      \tcount\ttotalCount\thash\tsize\tmodifyTime         \tname\n")
	}
}

var Cmd = &cobra.Command{
	Use:     "ls",
	Example: "kfs-cli ls .",
	Args:    cobra.RangeArgs(0, 1),
	Run:     runList,
}

func init() {
	Cmd.PersistentFlags().Bool(HumanizeStr, true, "")
}

func runList(cmd *cobra.Command, args []string) {
	var err error
	defer ExitWithError(err)
	serverType := viper.GetString(ServerTypeStr)
	serverAddr := viper.GetString(ServerAddrStr)
	branchName := viper.GetString(BranchNameStr)
	humanize := cmd.Flag(HumanizeStr).Value.String()
	fmt.Printf("%s=%s\n", ServerTypeStr, serverType)
	fmt.Printf("%s=%s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s=%s\n", BranchNameStr, branchName)
	p := ""
	if len(args) != 0 {
		p = args[0]
	}
	switch serverType {
	case ServerTypeLocal:
		err = local(serverAddr, branchName, p, humanize)
	case ServerTypeRemote:
		err = remote(serverAddr, branchName, p, humanize)
	default:
		err = InvalidServerType
	}
}
