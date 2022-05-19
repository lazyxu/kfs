package list

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	core "github.com/lazyxu/kfs/core/local"
)

func local(addr string, branchName string, p string, human string) error {
	kfsCore, _, err := core.New(addr)
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	ctx := context.Background()
	dirItems, err := kfsCore.List(ctx, branchName, p)
	if err != nil {
		return err
	}
	printHeader(len(dirItems))
	for _, dirItem := range dirItems {
		modifyTime := time.Unix(0, int64(dirItem.ModifyTime)).Format("2006-01-02 15:04:05")
		fmt.Printf("%s\t%5d\t%10d\t%s\t%s\t%s\t%s\n",
			os.FileMode(dirItem.Mode).String(), dirItem.Count, dirItem.TotalCount, dirItem.Hash[:4],
			humanize.Bytes(dirItem.Size), modifyTime, dirItem.Name)
	}
	return nil
}
