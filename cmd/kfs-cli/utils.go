package main

import (
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
)

func getFS(serverType string, serverAddr string) (core.FS, error) {
	switch serverType {
	case ServerTypeRemote:
		return client.New(serverAddr), nil
	}
	return nil, InvalidServerType
}
