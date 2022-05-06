package main

import (
	"context"
	"fmt"
	"os"
	"time"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func main() {
	if err := test(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func test() error {
	dbFileName := "kfs.db"
	db, err := sqlite.Open(dbFileName)
	if err != nil {
		return err
	}

	err = db.Reset()
	if err != nil {
		return err
	}

	ctx := context.Background()

	file1 := sqlite.NewFileFromBytes([]byte(nil), "")
	file2 := sqlite.NewFileFromBytes([]byte("abc"), "txt")
	err = db.WriteFile(ctx, file1)
	if err != nil {
		return err
	}

	err = db.WriteFile(ctx, file2)
	if err != nil {
		return err
	}

	now := uint64(time.Now().Nanosecond())
	dir, err := db.WriteDir(ctx, []sqlite.DirItem{
		sqlite.NewDirItem(file1, "emptyFile", 0o700, now, now, now, now, ""),
		sqlite.NewDirItem(file2, "aaa.txt", 0o555, now, now, now, now, ""),
	})
	if err != nil {
		return err
	}
	root, err := db.WriteDir(ctx, []sqlite.DirItem{
		sqlite.NewDirItem(dir, "data", 0o777, now, now, now, now, ""),
		sqlite.NewDirItem(file2, "bbb.txt", 0o555, now, now, now, now, ""),
	})
	if err != nil {
		return err
	}

	err = db.WriteBranch(ctx, sqlite.NewBranch("default", "no description", root))
	if err != nil {
		return err
	}

	count, err := db.FileCount(ctx)
	if err != nil {
		return err
	}
	println("FileCount", count)

	count, err = db.DirCount(ctx)
	if err != nil {
		return err
	}
	println("DirCount", count)

	count, err = db.DirItemCount(ctx)
	if err != nil {
		return err
	}
	println("DirItemCount", count)

	count, err = db.BranchCount(ctx)
	if err != nil {
		return err
	}
	println("BranchCount", count)

	if err = db.Close(); err != nil {
		return err
	}

	fi, err := os.Stat(dbFileName)
	if err != nil {
		return err
	}

	fmt.Printf("%s size: %v\n", dbFileName, fi.Size())
	return nil
}
