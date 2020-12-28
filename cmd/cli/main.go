package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "kfs",
	Short: "Kfs is a distributed file system named koala file system.",
	Long:  `Kfs is a file system designed to manage your files as easily as possible.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := initStorage()
		go initFuse(s)
		initGrpc(s)
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var root string

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&root, "kfs-root", "", "stores config file and objects")

	RootCmd.PersistentFlags().StringP("storage", "s", "", "storage type")
	viper.BindPFlag("storage", RootCmd.Flag("storage"))

	RootCmd.PersistentFlags().String("fuse-lib", "", "fuse library")
	viper.BindPFlag("fuse-lib", RootCmd.Flag("fuse-lib"))
	RootCmd.PersistentFlags().String("fuse-mount-point", "", "mount point for fuse")
	viper.BindPFlag("fuse-mount-point", RootCmd.Flag("fuse-mount-point"))

	RootCmd.PersistentFlags().Int("grpc-web-http-port", 9091, "http port")
	viper.BindPFlag("grpc-web-http-port", RootCmd.Flag("grpc-web-http-port"))
}

func initConfig() {
	if root != "" {
		viper.AddConfigPath(root)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		root = path.Join(home, ".kfs")
		viper.AddConfigPath(root)
	}
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can not read config:", viper.ConfigFileUsed())
	}
	viper.Set("kfs-root", root)
}
