package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

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
	viper.SetConfigName(".kfs-server")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	configFile := path.Join(home, ".kfs-server.json")
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
}
