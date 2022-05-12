package main

import (
	"context"
	"fmt"
	"os"

	core "github.com/lazyxu/kfs/core/local"
)

func main() {
	if err := test(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func test() error {
	//os.RemoveAll("tmp")
	// new
	// delete
	// move
	// cp
	kfs, exist, err := core.New("tmp")
	if err != nil {
		return err
	}
	defer kfs.Close()
	ctx := context.Background()
	branchName := "default"
	if !exist {
		err = kfs.Backup(ctx, "../..", branchName)
		if err != nil {
			return err
		}
	}
	dirItems, err := kfs.Ls(ctx, branchName, ".git")
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		fmt.Printf("%+v\n", dirItem.Name)
	}
	err = kfs.Delete(ctx, branchName, ".git", "refs")
	if err != nil {
		return err
	}
	println("------delete /.git/refs")
	dirItems, err = kfs.Ls(ctx, branchName, ".git")
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		fmt.Printf("%+v\n", dirItem.Name)
	}
	return nil
}
