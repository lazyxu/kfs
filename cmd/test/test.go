package main

import (
	"context"
	"fmt"
	"os"

	backup "github.com/lazyxu/kfs/backup/local"
)

func main() {
	if err := test(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func test() error {
	//os.RemoveAll("tmp")
	b, err := backup.New("tmp")
	if err != nil {
		return err
	}
	defer b.Close()
	ctx := context.Background()
	return b.Upload(ctx, "../..", "default")
}
