package main

import (
	"io"

	"github.com/spf13/cobra"
)

func catCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "cat",
		Example: "kfs-cli cat test.txt",
		Args:    cobra.RangeArgs(1, 1),
		Run:     runCat,
	}
}

func runCat(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(cmd, err)
	}()

	fs, branchName, _ := loadFs(cmd)

	srcPath := args[0]

	err = fs.Cat(cmd.Context(), branchName, srcPath, func(r io.Reader, size int64) error {
		_, e := io.CopyN(cmd.OutOrStdout(), r, size)
		if e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		return
	}
}
