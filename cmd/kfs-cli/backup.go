package main

import (
	"context"

	"github.com/lazyxu/kfs/cmd/kfs-cli/utils"

	"github.com/spf13/viper"

	core "github.com/lazyxu/kfs/core/local"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "backup files",
	Example: "kfs backup .",
	Run:     runBackup,
}

func runBackup(cmd *cobra.Command, args []string) {
	backupPath := ""
	if len(args) != 0 {
		backupPath = args[0]
	}
	kfsCore, _, err := core.New(viper.GetString(utils.ServerAddrStr))
	if err != nil {
		panic(err)
	}
	defer kfsCore.Close()
	ctx := context.Background()
	err = kfsCore.Backup(ctx, backupPath, viper.GetString(utils.BranchNameStr))
	if err != nil {
		panic(err)
	}
}
