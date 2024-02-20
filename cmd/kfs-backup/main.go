package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	err := rootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var verbose bool
var scanOnly bool
var srcPath string
var ignores []string

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "kfs-electron",
		Args: cobra.RangeArgs(0, 0),
		Run:  runRoot,
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	cmd.PersistentFlags().BoolVar(&scanOnly, "scan-only", false, "only scan files, not upload")
	cmd.PersistentFlags().StringVar(&srcPath, "src", "", "src path")
	cmd.PersistentFlags().StringSliceVar(&ignores, "ignore", []string{}, "ignores")
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) {
	doScan(cmd.Context(), srcPath, ignores, verbose)
}
