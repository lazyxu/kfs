package main

import (
	"github.com/spf13/cobra"
)

func resetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Example: "kfs-cli reset",
		Args:    cobra.RangeArgs(0, 0),
		Run:     runReset,
	}
}

func runReset(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

	err = fs.Reset(cmd.Context(), branchName)
}
