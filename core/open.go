package core

import (
	"context"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
	"os"
)

func (fs *KFS) Open(ctx context.Context, branchName string, filePath string) (mode os.FileMode, rc storage.SizedReadCloser, dirItems []sqlite.DirItem, err error) {
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
