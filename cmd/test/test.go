package main

import (
	"context"
	"fmt"
	"os"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

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
		err = kfsCore.Upload(ctx, branchName, "", "../..")
		if err != nil {
			return err
		}
	}
	err = kfsCore.List(ctx, branchName, ".git", nil, func(item sqlite.IDirItem) error {
		fmt.Printf("%+v\n", item.GetName())
		return nil
	})
	if err != nil {
		return err
	}
	_, _, err = kfsCore.Remove(ctx, branchName, ".git", "refs")
	if err != nil {
		return err
	}
	println("------delete /.git/refs")
	err = kfsCore.List(ctx, branchName, ".git", nil, func(item sqlite.IDirItem) error {
		fmt.Printf("%+v\n", item.GetName())
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
