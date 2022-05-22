package upload

import (
	"fmt"

	. "github.com/lazyxu/kfs/cmd/kfs-cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "upload",
	Example: "kfs-cli upload -p path filePath",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runUpload,
}

func init() {
	Cmd.PersistentFlags().StringP(PathStr, "p", "", "")
	Cmd.PersistentFlags().StringP(ChunkSizeStr, "b", "1 MiB", "[1 KiB, 1 GiB]")
}

const fileChunkSize = 1024 * 1024

func runUpload(cmd *cobra.Command, args []string) {
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

	// TODO: SET chunk bytes.
	//fileChunkSize := cmd.Flag(ChunkSizeStr)
	//humanize.ParseBytes()
	p := cmd.Flag(PathStr).Value.String()
	filename := args[0]

	switch serverType {
	case ServerTypeLocal:
		// TODO
	case ServerTypeRemote:
		err = remote(cmd.Context(), serverAddr, filename, branchName, p)
	default:
		err = InvalidServerType
	}
}