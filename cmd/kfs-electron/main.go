package main

import (
	"encoding/json"
	"fmt"
	"github.com/lazyxu/kfs/cmd/kfs-electron/db/gosqlite"
	"github.com/spf13/cobra"
	"net"
	"os"
)

func main() {
	err := rootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const PortStr = "port"

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "kfs-electron",
		Args: cobra.RangeArgs(0, 0),
		Run:  runRoot,
	}
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	cmd.PersistentFlags().StringP(PortStr, "p", "0", "local web server port")
	return cmd
}

var db *gosqlite.DB

func runRoot(cmd *cobra.Command, args []string) {
	var err error
	db, err = gosqlite.NewDb("electron.db")
	if err != nil {
		panic(err)
	}

	portStr := cmd.Flag(PortStr).Value.String()
	lis, err := net.Listen("tcp", "0.0.0.0:"+portStr)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(os.Stdout, "KFS electron web server listening at:", lis.Addr().String())
	if err != nil {
		panic(err)
	}
	if err != nil {
		return
	}
	if os.Getenv("KFS_CONFIG_PATH") != "" {
		filePath := os.Getenv("KFS_CONFIG_PATH")
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			panic(err)
		}
		m["port"] = lis.Addr().(*net.TCPAddr).Port
		data, err = json.Marshal(m)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(filePath, data, 0o600)
		if err != nil {
			panic(err)
		}
	}
	webServer(lis)
}
