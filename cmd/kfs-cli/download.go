package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func downloadCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "download",
		Example: "kfs-cli download -p path filePath",
		Args:    cobra.RangeArgs(0, 1),
		Run:     runDownload,
	}
	cmd.PersistentFlags().StringP(PathStr, "p", "", "override the path")
	cmd.PersistentFlags().String(DirPathStr, "", "move into dir")
	cmd.PersistentFlags().BoolP(VerboseStr, "v", false, "verbose")
	cmd.PersistentFlags().IntP(ConcurrentStr, "c", 1, "concurrent")
	cmd.PersistentFlags().StringP(EncoderStr, "e", "", "[\"\", \"lz4\"]")
	cmd.PersistentFlags().Bool(CpuProfilerStr, false, "cpu profile")
	cmd.PersistentFlags().StringP(ChunkSizeStr, "b", "1 MiB", "[1 KiB, 1 GiB]")
	return cmd
}

func runDownload(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

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
	var srcPath string
	if len(args) > 0 {
		srcPath = args[0]
	}
	if dstPath == "" {
		dstPath = srcPath
	}
	if dstPath == "" {
		err = fmt.Errorf("unknown dstPath")
	}

	//var uploadProcess core.UploadProcess = &core.EmptyUploadProcess{}
	//if verbose {
	//	uploadProcess = &UploadProcessBar{}
	//} else {
	//	uploadProcess = &core.EmptyUploadProcess{}
	//}

	if cpuProfile {

	}

	filePath, err := fs.Download(cmd.Context(), branchName, dstPath, srcPath)
	if err != nil {
		return
	}
	cmd.Printf("Saving to '%s'\n", filePath)
}
