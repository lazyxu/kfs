package main

import (
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:     "reset",
	Example: "kfs-cli reset",
	Args:    cobra.RangeArgs(0, 0),
	Run:     runReset,
}

func runReset(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(err)
	}()

	fs, branchName := loadFs(cmd)

	err = fs.Reset(cmd.Context(), branchName)
}
