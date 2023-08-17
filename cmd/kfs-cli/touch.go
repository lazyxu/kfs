package main

import (
	"github.com/dustin/go-humanize"

	"github.com/spf13/cobra"
)

func touchCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "touch",
		Example: "kfs-cli touch filePath",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runTouch,
	}
	return cmd
}

func runTouch(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

	filePath := "/" + args[0]

	branch, commit, err := fs.Touch(cmd.Context(), branchName, filePath)
	if err != nil {
		return
	}
	cmd.Printf("hash=%s, commitId=%d, size=%s, count=%d\n", branch.Hash[:4], commit.CommitId, humanize.IBytes(commit.Size), commit.Count)
}
