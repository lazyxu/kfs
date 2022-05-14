package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"

	core "github.com/lazyxu/kfs/core/local"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Short:   "ls list files",
	Example: "kfs backup .",
	Run:     runList,
}

const (
	kfsRootStr    = "kfs-root"
	backupPathStr = "backup-path"
	branchNameStr = "branch-name"
	pathStr       = "path"
)

func runList(cmd *cobra.Command, args []string) {
	kfsRoot := viper.GetString(kfsRootStr)
	branchName := viper.GetString(branchNameStr)
	p := ""
	if len(args) != 0 {
		p = args[0]
	}
	kfsCore, _, err := core.New(kfsRoot)
	if err != nil {
		panic(err)
	}
	defer kfsCore.Close()
	ctx := context.Background()
	dirItems, err := kfsCore.List(ctx, branchName, formatPath(p)...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("kfsRoot=%s\n", kfsRoot)
	fmt.Printf("branch=%s\n", branchName)
	fmt.Printf("total %d\n", len(dirItems))
	if len(dirItems) != 0 {
		fmt.Printf("mode      \tcount\ttotalCount\thash\tsize\tmodifyTime         \tname\n")
	}
	for _, dirItem := range dirItems {
		modifyTime := time.Unix(0, int64(dirItem.ModifyTime)).Format("2006-01-02 15:04:05")
		fmt.Printf("%s\t%5d\t%10d\t%s\t%s\t%s\t%s\n",
			os.FileMode(dirItem.Mode).String(), dirItem.Count, dirItem.TotalCount, dirItem.Hash[:4],
			humanize.Bytes(dirItem.Size), modifyTime, dirItem.Name)
	}
}

func calMode(mode uint64) string {
	arr := make([]uint8, 10)
	for i := range arr {
		arr[i] = '-'
	}
	template := "drwxrwxrwx"
	for i := 0; i < len(arr); i++ {
		if mode&(1<<i) != 0 {
			arr[i] = template[i]
		}
	}
	return string(arr)
}
