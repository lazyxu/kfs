package upload

import (
	"context"

	"github.com/lazyxu/kfs/core"
)

func local(ctx context.Context, addr string, branchName string, srcPath string, dstPath string) error {
	kfsCore, _, err := core.New(addr)
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	// TODO: dstPath
	return kfsCore.Backup(ctx, srcPath, branchName)
}
