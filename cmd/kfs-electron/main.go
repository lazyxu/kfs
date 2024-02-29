package main

import (
	"encoding/json"
	"fmt"
	"github.com/lazyxu/kfs/rpc/client/local_file"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	println("main")
	err := rootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var portStr string
var configPath string

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "kfs-electron",
		Args: cobra.RangeArgs(0, 0),
		Run:  runRoot,
	}
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	cmd.PersistentFlags().StringVarP(&portStr, "port", "p", "0", "local web server port")
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "~/.kfs-config.json", "config path")
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) {
	fmt.Printf("runRoot: %+v\n", args)

	println("portStr", portStr)
	lis, err := net.Listen("tcp", "0.0.0.0:"+portStr)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(os.Stdout, "KFS electron web server listening at:", lis.Addr().String())
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	m["port"] = lis.Addr().(*net.TCPAddr).Port
	data, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, data, 0o600)
	if err != nil {
		panic(err)
	}
	local_file.NewDeviceIfNeeded(configPath)
	webServer(lis)
}
