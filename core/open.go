package core

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (fs *KFS) Open(ctx context.Context, branchName string, filePath string) (mode os.FileMode, rc dao.SizedReadCloser, dirItems []dao.DirItem, err error) {
	var hash string
	hash, mode, dirItems, err = fs.Db.Open(ctx, branchName, FormatPath(filePath))
	if err != nil {
		return
	}
	if mode.IsRegular() {
		rc, err = fs.S.ReadWithSize(hash)
	}
	return
}

func (fs *KFS) Open2(ctx context.Context, branchName string, filePath string, maxContentSize int64) (dirItemOpened dao.DirItemOpened, err error) {
	dirItemOpened.DirItem, dirItemOpened.DirItems, err = fs.Db.Open2(ctx, branchName, FormatPath(filePath))
	if err != nil {
		return
	}
	if os.FileMode(dirItemOpened.DirItem.Mode).IsRegular() {
		if dirItemOpened.DirItem.Size > uint64(maxContentSize) {
			dirItemOpened.ContentTooLarge = true
			return
		}
		var rc dao.SizedReadCloser
		rc, err = fs.S.ReadWithSize(dirItemOpened.Hash)
		if err != nil {
			return
		}
		defer rc.Close()
		buf := bytes.NewBuffer(nil)
		_, err = io.CopyN(buf, rc, rc.Size())
		if err != nil {
			return
		}
		dirItemOpened.Content = buf.Bytes()
	}
	return
}

func (fs *KFS) OpenFile(ctx context.Context, driverId uint64, filePath []string, maxContentSize int64) (rc dao.SizedReadCloser, tooLarge bool, err error) {
	f, err := fs.Db.GetDriverFile(ctx, driverId, filePath)
	if err != nil {
		return
	}
	if f.Size > uint64(maxContentSize) {
		tooLarge = true
		return
	}
	rc, err = fs.S.ReadWithSize(f.Hash)
	return
}
