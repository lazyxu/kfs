package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kfs",
	Short: "Kfs is file system used to backup files.",
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().String(ConfigFileStr, filepath.Join(home, ".kfs.json"), "the path for the config file")
	rootCmd.PersistentFlags().BoolP(VerboseStr, "v", false, "verbose")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(catCmd)
}
