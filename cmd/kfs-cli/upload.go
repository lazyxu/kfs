package main

import (
	"fmt"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:     "upload",
	Example: "kfs-cli upload -p path filePath",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runUpload,
}

func init() {
	uploadCmd.PersistentFlags().StringP(PathStr, "p", "", "override the path")
	uploadCmd.PersistentFlags().String(DirPathStr, "", "move into dir")
	uploadCmd.PersistentFlags().BoolP(VerboseStr, "v", false, "verbose")
	uploadCmd.PersistentFlags().StringP(ChunkSizeStr, "b", "1 MiB", "[1 KiB, 1 GiB]")
}

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
	dstPath := cmd.Flag(PathStr).Value.String()
	flag := cmd.Flag(VerboseStr).Value.String()
	verbose := flag != "false"
	srcPath := args[0]

	err = withFS(serverType, serverAddr, func(fs core.FS) error {
		var uploadProcess core.UploadProcess
		if verbose {
			uploadProcess = &UploadProcessBar{}
		} else {
			uploadProcess = &core.EmptyUploadProcess{}
		}
		branch, commit, err := fs.Upload(cmd.Context(), branchName, dstPath, srcPath, uploadProcess)
		if err != nil {
			return err
		}
		fmt.Printf("hash=%s, commitId=%d, size=%d, count=%d\n", branch.Hash[:4], commit.CommitId, commit.Size, commit.Count)
		return nil
	})
}
