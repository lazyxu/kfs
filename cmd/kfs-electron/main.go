package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	"github.com/lazyxu/kfs/rpc/client/local_file"

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
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".kfs-config.json", "config path")
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
	m := map[string]interface{}{}
	println("configPath", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("readfile error %v\n", err)
	} else {
		err = json.Unmarshal(data, &m)
		if err != nil {
			fmt.Printf("unmarshal error %v\n", err)
		}
	}
	m["port"] = lis.Addr().(*net.TCPAddr).Port
	data, err = json.MarshalIndent(m, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, data, 0o600)
	if err != nil {
		fmt.Printf("WriteFile error %v\n", err)
		panic(err)
	}

	serverAddr := m["socketServer"].(string)
	deviceID := local_file.NewDeviceIfNeeded(configPath)
	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}
	err = fs.UploadDir(context.TODO(), deviceID, 4, "/", "/home/xuliang/kfs/cmd/kfs-electron/task/local_file_filter", core.UploadDirConfig{
		UploadDirProcess: &core.DefaultUploadDirProcessor{},
		Encoder:          "",
		Concurrent:       1,
		Verbose:          false,
	})
	if err != nil {
		fmt.Printf("WriteFile error %v\n", err)
		panic(err)
	}

	webServer(lis)
}
