package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "init config",
	Example: "kfs init ./tmp -b master",
	Run:     runInit,
}

func init() {
	initCmd.PersistentFlags().StringP(branchNameStr, "b", "master", "")
}

func runInit(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, errors.New("expected kfs root dir"))
		os.Exit(1)
	}
	kfsRoot, err := filepath.Abs(args[0])
	if err != nil {
		panic(err)
	}
	viper.Set(kfsRootStr, kfsRoot)
	branchName := cmd.Flag(branchNameStr).Value.String()
	viper.Set(branchNameStr, branchName)
	err = viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}
