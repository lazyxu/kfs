package main

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/lazyxu/kfs/rpc/client/local_file"
	"os"
	"path/filepath"

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

	cmd.PersistentFlags().Uint64VarP(&driverId, "driver", "d", invalidDriverId, "driver id")

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", filepath.Join(home, "kfs-config.json"), "config path")
	return cmd
}

func initConfigFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(path)
			if err != nil {
				fmt.Printf("创建配置文件失败： %s\n", path)
				os.Exit(1)
			}
			defer f.Close()
			_, err = f.WriteString("{}")
			if err != nil {
				fmt.Printf("初始化配置文件失败： %s\n", path)
				os.Exit(1)
			}
		} else {
			fmt.Printf("读取配置文件失败： %s\n", path)
			os.Exit(1)
		}
	}
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
		initConfigFile(configPath)
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("读取配置文件失败： %s\n", configPath)
			return
		}
		socketServer := json.ConfigCompatibleWithStandardLibrary.Get(data, "socketServer").ToString()
		if socketServer == "" {
			fmt.Printf("请配置socketServer\n")
			os.Exit(1)
		}
		deviceId := local_file.NewDeviceIfNeeded(configPath)
		doUpload(ctx, deviceId, socketServer, driverId, srcPath, ignores, verbose)
	}
}
