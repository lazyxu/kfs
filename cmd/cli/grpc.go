package main

import (
	"fmt"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/lazyxu/kfs/warpper/grpcweb"

	"github.com/spf13/viper"
)

func initGrpc(s storage.Storage) {
	httpPort := viper.GetInt("grpc-web-http-port")
	fmt.Println("grpc", httpPort)
	if s == nil {
		panic("storage is nil")
	}
	grpcweb.Start(httpPort, s)
}
