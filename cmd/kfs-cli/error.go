package main

import (
	"github.com/spf13/cobra"
)

func ExitWithError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err)
	}
}
