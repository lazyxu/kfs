package main

import (
	"fmt"

	"github.com/lazyxu/kfs/rpc/client"
	"github.com/spf13/viper"
)

func loadFs() (*client.RpcFs, string) {
	grpcServerAddr := viper.GetString(GrpcServerAddrStr)
	socketServerAddr := viper.GetString(SocketServerAddrStr)
	branchName := viper.GetString(BranchNameStr)
	fmt.Printf("%s: %s\n", BranchNameStr, branchName)
	return &client.RpcFs{
		GrpcServerAddr:   grpcServerAddr,
		SocketServerAddr: socketServerAddr,
	}, branchName
}
