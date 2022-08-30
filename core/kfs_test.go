package core

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lazyxu/kfs/db/gosqlite"
	"github.com/lazyxu/kfs/db/mysql"

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
		return New(gosqlite.FuncNew(testRootDir), storage.FuncNew(testRootDir, storage.NewStorage0))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage1Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(gosqlite.FuncNew(testRootDir), storage.FuncNew(testRootDir, storage.NewStorage1))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage2Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(gosqlite.FuncNew(testRootDir), storage.FuncNew(testRootDir, storage.NewStorage2))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage3Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(gosqlite.FuncNew(testRootDir), storage.FuncNew(testRootDir, storage.NewStorage3))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(gosqlite.FuncNew(testRootDir), storage.FuncNew(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

var dataSourceName = "root:12345678@/kfs?parseTime=true&multiStatements=true"

func BenchmarkMysqlStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(mysql.FuncNew(dataSourceName), storage.FuncNew(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage5Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(mysql.FuncNew(dataSourceName), storage.FuncNew(testRootDir, storage.NewStorage5))
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
		err = kfsCore.ResetBranch(ctx, branchName)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		wg := sync.WaitGroup{}
		wg.Add(fileCount)
		for j := 0; j < fileCount; j++ {
			//go func(j int) {
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
			go func() {
				defer wg.Done()
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
			}()
			//}(j)
		}
		wg.Wait()
		b.StopTimer()
	}
}
