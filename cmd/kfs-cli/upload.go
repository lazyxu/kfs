package main

import (
	"strconv"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
)

func uploadCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "upload",
		Example: "kfs-cli upload -p path filePath",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runUpload,
	}

	cmd.PersistentFlags().StringP(PathStr, "p", "", "override the path")
	cmd.PersistentFlags().String(DirPathStr, "", "move into dir")
	cmd.PersistentFlags().IntP(ConcurrentStr, "c", 1, "concurrent")
	cmd.PersistentFlags().StringP(EncoderStr, "e", "", "[\"\", \"lz4\"]")
	cmd.PersistentFlags().Bool(CpuProfilerStr, false, "cpu profile")
	cmd.PersistentFlags().StringP(ChunkSizeStr, "b", "1 MiB", "[1 KiB, 1 GiB]")
	return cmd
}

func runUpload(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, verbose := loadFs(cmd)

	// TODO: SET chunk bytes.
	//fileChunkSize := cmd.Flag(ChunkSizeStr)
	//humanize.ParseBytes()
	dstPath := cmd.Flag(PathStr).Value.String()
	encoder := cmd.Flag(EncoderStr).Value.String()
	cpuProfile := cmd.Flag(CpuProfilerStr).Value.String() != "false"
	concurrent, err := strconv.Atoi(cmd.Flag(ConcurrentStr).Value.String())
	if err != nil {
		return
	}
	srcPath := args[0]

	var uploadProcess core.UploadProcess = &core.EmptyUploadProcess{}

	if cpuProfile {

	}

	commit, branch, err := fs.Upload(cmd.Context(), branchName, dstPath, srcPath, core.UploadConfig{
		Encoder:       encoder,
		UploadProcess: uploadProcess,
		Concurrent:    concurrent,
		Verbose:       verbose,
	})
	if err != nil {
		return
	}
	cmd.Printf("hash=%s, commitId=%d, size=%s, count=%d\n", commit.Hash[:4], branch.CommitId, humanize.Bytes(branch.Size), branch.Count)
}
