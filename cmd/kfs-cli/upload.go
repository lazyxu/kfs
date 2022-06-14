package main

import (
	"fmt"
	"strconv"

	"github.com/dustin/go-humanize"

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
	uploadCmd.PersistentFlags().IntP(ConcurrentStr, "c", 1, "concurrent")
	uploadCmd.PersistentFlags().StringP(EncoderStr, "e", "", "[\"\", \"lz4\"]")
	uploadCmd.PersistentFlags().Bool(CpuProfilerStr, false, "cpu profile")
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
	verbose := cmd.Flag(VerboseStr).Value.String() != "false"
	encoder := cmd.Flag(EncoderStr).Value.String()
	cpuProfile := cmd.Flag(CpuProfilerStr).Value.String() != "false"
	concurrent, err := strconv.Atoi(cmd.Flag(ConcurrentStr).Value.String())
	if err != nil {
		return
	}
	srcPath := args[0]

	var uploadProcess core.UploadProcess = &core.EmptyUploadProcess{}
	//if verbose {
	//	uploadProcess = &UploadProcessBar{}
	//} else {
	//	uploadProcess = &core.EmptyUploadProcess{}
	//}

	if cpuProfile {

	}

	fs, err := getFS(serverType, serverAddr)
	if err != nil {
		return
	}
	defer fs.Close()

	branch, commit, err := fs.Upload(cmd.Context(), branchName, dstPath, srcPath, core.UploadConfig{
		Encoder:       encoder,
		UploadProcess: uploadProcess,
		Concurrent:    concurrent,
		Verbose:       verbose,
	})
	if err != nil {
		return
	}
	fmt.Printf("hash=%s, commitId=%d, size=%s, count=%d\n", branch.Hash[:4], commit.CommitId, humanize.Bytes(commit.Size), commit.Count)
}
