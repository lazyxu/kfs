package list

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

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

func formatCount(mode uint64, count uint64) string {
	if !os.FileMode(mode).IsDir() {
		return strings.Repeat(" ", 5)
	}
	return fmt.Sprintf("%5d", count)
}

func printBody(dirItem sqlite.IDirItem, isHumanize bool) {
	modifyTime := time.Unix(0, int64(dirItem.GetModifyTime())).Format("2006-01-02 15:04:05")
	if isHumanize {
		fmt.Printf("%s\t%s\t     %s\t%s\t%s\t%s\t%s\n",
			os.FileMode(dirItem.GetMode()).String(),
			formatCount(dirItem.GetMode(), dirItem.GetCount()), formatCount(dirItem.GetMode(), dirItem.GetTotalCount()),
			dirItem.GetHash()[:4], humanize.Bytes(dirItem.GetSize()), modifyTime, dirItem.GetName())
	} else {
		fmt.Printf("%s\t%s\t     %s\t%s\t%d\t%s\t%s\n",
			os.FileMode(dirItem.GetMode()).String(),
			formatCount(dirItem.GetMode(), dirItem.GetCount()), formatCount(dirItem.GetMode(), dirItem.GetTotalCount()),
			dirItem.GetHash()[:4], dirItem.GetSize(), modifyTime, dirItem.GetName())
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
	defer func() {
		ExitWithError(err)
	}()
	serverType := viper.GetString(ServerTypeStr)
	serverAddr := viper.GetString(ServerAddrStr)
	branchName := viper.GetString(BranchNameStr)
	fmt.Printf("%s: %s\n", ServerTypeStr, serverType)
	fmt.Printf("%s: %s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s: %s\n", BranchNameStr, branchName)

	humanize := cmd.Flag(HumanizeStr).Value.String()
	p := ""
	if len(args) != 0 {
		p = args[0]
	}
	isHumanize := humanize == "true"

	switch serverType {
	case ServerTypeLocal:
		err = local(cmd.Context(), serverAddr, branchName, p, isHumanize)
	case ServerTypeRemote:
		err = remote(cmd.Context(), serverAddr, branchName, p, isHumanize)
	default:
		err = InvalidServerType
	}
}
