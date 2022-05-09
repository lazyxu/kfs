package main

import (
	"context"
	"fmt"
	"os"

	backup "github.com/lazyxu/kfs/backup/local"
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
	os.Remove(dbFileName)
	db, err := sqlite.Open(dbFileName)
	if err != nil {
		return err
	}

	err = db.Reset()
	if err != nil {
		return err
	}

	b, err := backup.New(dbFileName)
	if err != nil {
		return err
	}
	defer b.Close()
	ctx := context.Background()
	return b.Upload(ctx, "../..")
}
