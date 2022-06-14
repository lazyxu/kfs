package main

import (
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
)

func getFS(serverType string, serverAddr string) (core.FS, error) {
	switch serverType {
	case ServerTypeLocal:
		fs, _, err := core.New(serverAddr)
		if err != nil {
			return nil, err
		}
		return fs, nil
	case ServerTypeRemote:
		return client.New(serverAddr), nil
	}
	return nil, InvalidServerType
}
