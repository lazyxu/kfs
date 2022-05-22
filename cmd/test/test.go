package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lazyxu/kfs/core"
)

func main() {
	if err := test(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func test() error {
	//os.RemoveAll("tmp")
	kfsCore, exist, err := core.New("tmp")
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	ctx := context.Background()
	branchName := "default"
	if !exist {
		err = kfsCore.Backup(ctx, "../..", branchName)
		if err != nil {
			return err
		}
	}
	dirItems, err := kfsCore.List(ctx, branchName, ".git")
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		fmt.Printf("%+v\n", dirItem.Name)
	}
	_, err = kfsCore.Remove(ctx, branchName, ".git", "refs")
	if err != nil {
		return err
	}
	println("------delete /.git/refs")
	dirItems, err = kfsCore.List(ctx, branchName, ".git")
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		fmt.Printf("%+v\n", dirItem.Name)
	}
	return nil
}
