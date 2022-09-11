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

	"github.com/lazyxu/kfs/db/cgosqlite"
	"github.com/lazyxu/kfs/db/gosqlite"
	"github.com/lazyxu/kfs/db/mysql"

	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	storage "github.com/lazyxu/kfs/storage/local"
)

const (
	testRootDir      = "tmp"
	sqliteDataSource = "kfs.db"
)

func BenchmarkStorage0Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage0))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage1Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage1))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage2Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage2))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage3Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage3))
	}, branchName, fileCount, fileSize)
}

func BenchmarkStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

func BenchmarkCgoSqliteStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, cgosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

func BenchmarkCgoSqliteStorage4Upload10000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, cgosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

var mysqlDataSourceName = "root:12345678@/kfs?parseTime=true&multiStatements=true"

func BenchmarkMysqlStorage1Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage1))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage4Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage4Upload10000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage4))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage5Upload1000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 1000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage5Upload10000Files1000(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFiles(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func storageUploadFiles(b *testing.B, newKFS func() (*KFS, error), branchName string, fileCount int, fileSize int) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		kfsCore, err := newKFS()
		if err != nil {
			b.Error(err)
			return
		}
		err = kfsCore.Reset()
		if err != nil {
			b.Error(err)
			return
		}
		_, err = kfsCore.Checkout(context.Background(), "master")
		if err != nil {
			b.Error(err)
			return
		}
		b.StartTimer()
		wg := sync.WaitGroup{}
		wg.Add(fileCount)
		for j := 0; j < fileCount; j++ {
			//go func(j int) {
			fileName := strconv.Itoa(j)
			hash, content := storage.NewContent(strconv.Itoa(j) + strings.Repeat("y", fileSize) + "\n")
			mode := uint64(os.FileMode(0o700))
			now := uint64(time.Now().UnixNano())
			exist, err := kfsCore.S.Write(hash, func(f io.Writer, hasher io.Writer) (e error) {
				w := io.MultiWriter(f, hasher)
				_, e = io.CopyN(w, bytes.NewBuffer(content), int64(len(content)))
				return rpcutil.UnexpectedIfError(e)
			})
			if err != nil {
				b.Error(err)
				return
			}
			if exist {
				b.Error("should not exist")
				return
			}
			go func() {
				defer wg.Done()
				_, _, err = kfsCore.Db.UpsertDirItem(context.TODO(), branchName, FormatPath(fileName), dao.DirItem{
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
		err = kfsCore.Close()
		if err != nil {
			b.Error(err)
			return
		}
		b.StopTimer()
	}
}

func BenchmarkCgoSqliteStorage5Upload10000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, cgosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkGoSqliteStorage5Upload10000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage5Upload10000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 10000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkCgoSqliteStorage5Upload100000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 100000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, cgosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkGoSqliteStorage5Upload100000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 100000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(sqliteDataSource, gosqlite.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func BenchmarkMysqlStorage5Upload100000Files1000Batch(b *testing.B) {
	branchName := "master"
	fileCount := 100000
	fileSize := 1000
	storageUploadFilesBatch(b, func() (*KFS, error) {
		return New(dao.DatabaseNewFunc(mysqlDataSourceName, mysql.New), dao.StorageNewFunc(testRootDir, storage.NewStorage5))
	}, branchName, fileCount, fileSize)
}

func storageUploadFilesBatch(b *testing.B, newKFS func() (*KFS, error), branchName string, fileCount int, fileSize int) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		kfsCore, err := newKFS()
		if err != nil {
			b.Error(err)
			return
		}
		err = kfsCore.Reset()
		if err != nil {
			b.Error(err)
			return
		}
		_, err = kfsCore.Checkout(context.Background(), "master")
		if err != nil {
			b.Error(err)
			return
		}
		b.StartTimer()
		dirItems := make([]dao.DirItem, fileCount)
		for j := 0; j < fileCount; j++ {
			//go func(j int) {
			fileName := strconv.Itoa(j)
			hash, content := storage.NewContent(strconv.Itoa(j) + strings.Repeat("y", fileSize) + "\n")
			mode := uint64(os.FileMode(0o700))
			now := uint64(time.Now().UnixNano())
			exist, err := kfsCore.S.Write(hash, func(f io.Writer, hasher io.Writer) (e error) {
				w := io.MultiWriter(f, hasher)
				_, e = io.CopyN(w, bytes.NewBuffer(content), int64(len(content)))
				return rpcutil.UnexpectedIfError(e)
			})
			if err != nil {
				b.Error(err)
				return
			}
			if exist {
				b.Error("should not exist")
				return
			}
			dirItems[j] = dao.DirItem{
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
			}
		}
		_, _, err = kfsCore.Db.UpsertDirItems(context.TODO(), branchName, []string{}, dirItems)
		if err != nil {
			b.Error(err)
			return
		}
		err = kfsCore.Close()
		if err != nil {
			b.Error(err)
			return
		}
		b.StopTimer()
	}
}
