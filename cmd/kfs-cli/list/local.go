package list

import (
	"context"

	"github.com/lazyxu/kfs/core"
)

func local(ctx context.Context, addr string, branchName string, p string, isHumanize bool) error {
	kfsCore, _, err := core.New(addr)
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	dirItems, err := kfsCore.List(ctx, branchName, p)
	if err != nil {
		return err
	}
	printHeader(len(dirItems))
	for _, dirItem := range dirItems {
		printBody(dirItem, isHumanize)
	}
	return nil
}
