package mysql

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/lazyxu/kfs/dao"

	storage "github.com/lazyxu/kfs/storage/local"
)

func TestSqlite(t *testing.T) {
	dbFileName := "root:12345678@/kfs?parseTime=true&multiStatements=true"
	db, err := Open(dbFileName)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	err = db.Remove()
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Create()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	s, err := storage.NewStorage0("tmp")
	if err != nil {
		t.Error(err)
		return
	}
	err = s.Remove()
	if err != nil {
		t.Error(err)
		return
	}
	err = s.Create()
	if err != nil {
		t.Error(err)
		return
	}

	hash1, content1 := storage.NewContent("")
	_, err = storage.Write(s, hash1, bytes.NewReader(content1))
	if err != nil {
		t.Error(err)
		return
	}

	hash2, content2 := storage.NewContent("abc")
	_, err = storage.Write(s, hash2, bytes.NewReader(content2))
	if err != nil {
		t.Error(err)
		return
	}

	file1 := dao.NewFileByBytes(content1)
	file2 := dao.NewFileByBytes(content2)
	err = db.WriteFile(ctx, file1)
	if err != nil {
		t.Error(err)
		return
	}

	err = db.WriteFile(ctx, file2)
	if err != nil {
		t.Error(err)
		return
	}

	now := uint64(time.Now().UnixNano())
	dir, err := db.WriteDir(ctx, []dao.DirItem{
		dao.NewDirItem(file1, "emptyFile", 0o700, now, now, now, now),
		dao.NewDirItem(file2, "aaa.txt", 0o555, now, now, now, now),
	})
	if err != nil {
		t.Error(err)
		return
	}
	root, err := db.WriteDir(ctx, []dao.DirItem{
		dao.NewDirItem(dir, "data", 0o777, now, now, now, now),
		dao.NewDirItem(file2, "bbb.txt", 0o555, now, now, now, now),
	})
	if err != nil {
		t.Error(err)
		return
	}

	branchName := "default"
	commit := dao.NewCommit(root, branchName, "")
	err = db.WriteCommit(ctx, &commit)
	if err != nil {
		t.Error(err)
		return
	}

	err = db.insertBranch(ctx, db.db, dao.NewBranch(branchName, commit, root))
	if err != nil {
		t.Error(err)
		return
	}

	count, err := db.FileCount(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if count != 2 {
		t.Errorf("invalid FileCount: expected %d, actual %d", 2, count)
	}

	count, err = db.DirCount(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if count != 2 {
		t.Errorf("invalid DirCount: expected %d, actual %d", 2, count)
	}

	count, err = db.DirItemCount(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if count != 4 {
		t.Errorf("invalid DirItemCount: expected %d, actual %d", 4, count)
	}

	count, err = db.BranchCount(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if count != 1 {
		t.Errorf("invalid BranchCount: expected %d, actual %d", 1, count)
	}
}
