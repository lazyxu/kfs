package main

import (
	"context"
	"fmt"
	"io"
	"os"

	core "github.com/lazyxu/kfs/core/local"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var catCmd = &cobra.Command{
	Use:     "cat",
	Short:   "cat file",
	Example: "kfs cat test.txt",
	Run:     runCat,
}

func runCat(cmd *cobra.Command, args []string) {
	kfsRoot := viper.GetString(kfsRootStr)
	branchName := viper.GetString(branchNameStr)
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "expected file path")
		return
	}
	p := args[0]
	kfsCore, _, err := core.New(kfsRoot)
	if err != nil {
		panic(err)
	}
	defer kfsCore.Close()
	ctx := context.Background()
	readerCloser, err := kfsCore.Cat(ctx, branchName, formatPath(p)...)
	if err != nil {
		panic(err)
	}
	defer readerCloser.Close()
	fmt.Printf("kfsRoot=%s\n", kfsRoot)
	fmt.Printf("branch=%s\n", branchName)
	_, err = io.Copy(os.Stdout, readerCloser)
	if err != nil {
		panic(err)
	}
}
