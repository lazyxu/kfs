package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExitWithError(cmd *cobra.Command, err error) {
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), err)
	}
}
