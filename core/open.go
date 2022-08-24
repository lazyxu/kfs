package core

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	storage "github.com/lazyxu/kfs/storage/local"
	"os"
)

func (fs *KFS) Open(ctx context.Context, branchName string, filePath string) (mode os.FileMode, rc storage.SizedReadCloser, dirItems []dao.DirItem, err error) {
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
