package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lazyxu/kfs/rpc/client"
	"github.com/spf13/viper"
)

func loadFs(cmd *cobra.Command) (*client.RpcFs, string, bool) {
	loadConfigFile(cmd)
	verbose := cmd.Flag(VerboseStr).Value.String() != "false"
	grpcServerAddr := viper.GetString(GrpcServerStr)
	serverServerAddr := viper.GetString(SocketServerStr)
	branchName := viper.GetString(BranchNameStr)
	if verbose {
		fmt.Printf("%s: %s\n", BranchNameStr, branchName)
	}
	return &client.RpcFs{
		GrpcServerAddr:   grpcServerAddr,
		SocketServerAddr: serverServerAddr,
	}, branchName, verbose
}

func loadConfigFile(cmd *cobra.Command) {
	configFilePath := cmd.Flag(ConfigFileStr).Value.String()
	configFile, err := filepath.Abs(configFilePath)
	if err != nil {
		panic(err)
	}
	extIndex := strings.LastIndexByte(configFile, '.')
	ext := configFile[extIndex+1:]
	fileName := configFile[:extIndex]
	pathSeparatorIndex := strings.LastIndexByte(fileName, os.PathSeparator)
	dir := configFile[:pathSeparatorIndex+1]
	fileName = configFile[pathSeparatorIndex+1:]
	viper.AddConfigPath(dir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(ext)
	viper.AutomaticEnv()
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
