package cat

import (
	"fmt"
	"io"
	"os"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"
	"github.com/lazyxu/kfs/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "cat",
	Example: "kfs-cli cat test.txt",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCat,
}

func runCat(cmd *cobra.Command, args []string) {
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

	p := args[0]

	var readerCloser io.ReadCloser
	switch serverType {
	case ServerTypeLocal:
		readerCloser, err = core.Cat(cmd.Context(), serverAddr, branchName, p)
	case ServerTypeRemote:
	default:
		err = InvalidServerType
	}

	if err != nil {
		return
	}
	defer readerCloser.Close()
	_, err = io.Copy(os.Stdout, readerCloser)
	if err != nil {
		return
	}
}
