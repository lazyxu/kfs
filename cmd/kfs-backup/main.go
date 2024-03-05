package main

import (
	"encoding/json"
	"fmt"
	"github.com/lazyxu/kfs/rpc/client/local_file"
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

var serverAddr string
var driverId uint64

var configPath string

const invalidDriverId uint64 = 18446744073709551615

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

	cmd.PersistentFlags().StringVar(&serverAddr, "server", "", "server address")
	cmd.PersistentFlags().Uint64VarP(&driverId, "driver", "d", invalidDriverId, "driver id")

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "~/.kfs-config.json", "config path")
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if scanOnly {
		doScan(ctx, srcPath, ignores, verbose)
	} else {
		if driverId == invalidDriverId {
			fmt.Printf("请输入正确的云盘ID\n")
			return
		}
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("读取配置文件失败： %s\n", configPath)
			return
		}
		deviceId := local_file.NewDeviceIfNeeded(configPath)
		m := map[string]interface{}{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			panic(err)
		}
		if serverAddr == "" {
			serverAddr = m["socketServer"].(string)
		}
		doUpload(ctx, deviceId, serverAddr, driverId, srcPath, ignores, verbose)
	}
}