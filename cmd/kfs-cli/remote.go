package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var remoteCmd = &cobra.Command{
	Use:     "remote",
	Short:   "set remote address",
	Example: "kfs-cli remote localhost:1123",
	Args:    cobra.RangeArgs(1, 1),
	Run:     runRemote,
}

func runRemote(cmd *cobra.Command, args []string) {
	remoteAddr := args[0]
	viper.Set(remoteAddrStr, remoteAddr)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("set remote addr to", remoteAddr)
}
