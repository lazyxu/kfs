package main

import (
	"fmt"
	"strconv"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCmd = &cobra.Command{
	Use:     "download",
	Example: "kfs-cli download -p path filePath",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runDownload,
}

func init() {
	downloadCmd.PersistentFlags().StringP(PathStr, "p", "", "override the path")
	downloadCmd.PersistentFlags().String(DirPathStr, "", "move into dir")
	downloadCmd.PersistentFlags().BoolP(VerboseStr, "v", false, "verbose")
	downloadCmd.PersistentFlags().IntP(ConcurrentStr, "c", 1, "concurrent")
	downloadCmd.PersistentFlags().StringP(EncoderStr, "e", "", "[\"\", \"lz4\"]")
	downloadCmd.PersistentFlags().Bool(CpuProfilerStr, false, "cpu profile")
	downloadCmd.PersistentFlags().StringP(ChunkSizeStr, "b", "1 MiB", "[1 KiB, 1 GiB]")
}

func runDownload(cmd *cobra.Command, args []string) {
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

	filePath, err := fs.Download(cmd.Context(), branchName, dstPath, srcPath, core.UploadConfig{
		Encoder:       encoder,
		UploadProcess: uploadProcess,
		Concurrent:    concurrent,
		Verbose:       verbose,
	})
	if err != nil {
		return
	}
	fmt.Printf("Saving to '%s'\n", filePath)
}
