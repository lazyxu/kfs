package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func main() {
	err := rootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "kfs",
		Short: "Kfs is file system used to backup files.",
	}
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	cmd.PersistentFlags().String(ConfigFileStr, filepath.Join(home, ".kfs.json"), "the path for the config file")
	cmd.PersistentFlags().BoolP(VerboseStr, "v", false, "verbose")
	cmd.AddCommand(initCmd())
	cmd.AddCommand(branchCmd())
	cmd.AddCommand(checkoutCmd())
	cmd.AddCommand(uploadCmd())
	cmd.AddCommand(downloadCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(resetCmd())
	cmd.AddCommand(catCmd())
	cmd.AddCommand(touchCmd())
	return cmd
}
