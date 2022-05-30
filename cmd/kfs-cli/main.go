package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kfs",
	Short: "Kfs is file system used to backup files.",
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".kfs")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	configFile := path.Join(home, ".kfs.json")
	_, err = os.Stat(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.Write([]byte("{}"))
		if err != nil {
			panic(err)
		}
	}
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Can not read config: %s\n", viper.ConfigFileUsed())
	}
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(catCmd)
}
