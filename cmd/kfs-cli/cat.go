package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lazyxu/kfs/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var catCmd = &cobra.Command{
	Use:     "cat",
	Example: "kfs-cli cat test.txt",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runCat,
}

func runCat(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		ExitWithError(err)
	}()
	serverAddr := viper.GetString(ServerAddrStr)
	branchName := viper.GetString(BranchNameStr)
	fmt.Printf("%s: %s\n", ServerAddrStr, serverAddr)
	fmt.Printf("%s: %s\n", BranchNameStr, branchName)

	p := args[0]

	var readerCloser io.ReadCloser
	readerCloser, err = core.Cat(cmd.Context(), serverAddr, branchName, p)

	if err != nil {
		return
	}
	defer readerCloser.Close()
	_, err = io.Copy(os.Stdout, readerCloser)
	if err != nil {
		return
	}
}
