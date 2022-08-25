package core

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	storage "github.com/lazyxu/kfs/storage/local"
)

const (
	testRootDir = "tmp"
)

func BenchmarkStorage0Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		kfs, _, err := NewWithSqlite(testRootDir, storage.NewStorage0)
		return kfs, err
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage1Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		kfs, _, err := NewWithSqlite(testRootDir, storage.NewStorage1)
		return kfs, err
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage2Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		kfs, _, err := NewWithSqlite(testRootDir, storage.NewStorage2)
		return kfs, err
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage3Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		kfs, _, err := NewWithSqlite(testRootDir, storage.NewStorage3)
		return kfs, err
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		kfs, _, err := NewWithSqlite(testRootDir, storage.NewStorage4)
		return kfs, err
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return NewWithMysql(testRootDir, storage.NewStorage4)
	}, branchName, fileCount, fileSize)
}

func storageUploadFiles(b *testing.B, newKFS func() (*KFS, error), branchName string, fileCount int, fileSize int) {
	kfsCore, err := newKFS()
	if err != nil {
		b.Error(err)
		return
	}
	defer kfsCore.Close()
	ctx := context.TODO()
	for i := 0; i < b.N; i++ {
		err = kfsCore.Reset(ctx, branchName)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		for j := 0; j < fileCount; j++ {
			fileName := strconv.Itoa(j)
			hash, content := storage.NewContent(strconv.Itoa(j) + strings.Repeat("y", fileSize) + "\n")
			mode := uint64(os.FileMode(0o700))
			now := uint64(time.Now().UnixNano())
			exist, err := kfsCore.S.WriteFn(hash, func(f io.Writer, hasher io.Writer) (e error) {
				w := io.MultiWriter(f, hasher)
				_, e = io.CopyN(w, bytes.NewBuffer(content), int64(len(content)))
				return rpcutil.UnexpectedIfError(e)
			})
			if exist {
				b.Error("should not exist")
				return
			}
			_, _, err = kfsCore.Db.UpsertDirItem(ctx, branchName, FormatPath(fileName), dao.DirItem{
				Hash:       hash,
				Name:       fileName,
				Mode:       mode,
				Size:       uint64(len(content)),
				Count:      1,
				TotalCount: 1,
				CreateTime: now,
				ModifyTime: now,
				ChangeTime: now,
				AccessTime: now,
			})
			if err != nil {
				b.Error(err)
				return
			}
		}
		b.StopTimer()
	}
}
