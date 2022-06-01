package main

import (
	"strings"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
)

func formatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	if splitPath[0] == "" {
		splitPath = splitPath[1:]
	}
	return splitPath
}

func withFS(serverType string, serverAddr string, fn func(fs core.FS) error) error {
	switch serverType {
	case ServerTypeLocal:
		fs, _, err := core.New(serverAddr)
		if err != nil {
			return err
		}
		defer fs.Close()
		return fn(fs)
	case ServerTypeRemote:
		fs := client.New(serverAddr)
		return fn(fs)
	default:
		return InvalidServerType
	}
}
