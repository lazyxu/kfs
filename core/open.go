package core

import (
	"context"
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
