package core

import (
	"context"
	"testing"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func TestKFS(t *testing.T) {
	kfsCore, _, err := New("tmp")
	if err != nil {
		t.Error(err)
		return
	}
	defer kfsCore.Close()
	ctx := context.Background()
	branchName := "default"
	_, _, err = kfsCore.Upload(ctx, branchName, "", ".", UploadConfig{
		UploadProcess: &EmptyUploadProcess{},
		Concurrent:    1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	count := 0
	err = kfsCore.List(ctx, branchName, "", nil, func(item sqlite.IDirItem) error {
		count++
		// fmt.Printf("%+v\n", item.GetName())
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, _, err = kfsCore.Remove(ctx, branchName, "kfs_test.go")
	if err != nil {
		t.Error(err)
		return
	}
	//println("------delete /kfs_test.go")
	err = kfsCore.List(ctx, branchName, "", nil, func(item sqlite.IDirItem) error {
		count--
		// fmt.Printf("%+v\n", item.GetName())
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if count != 1 {
		t.Errorf("invalid count: expected %d, actual %d", 1, count)
	}
}
