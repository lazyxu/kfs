package core

import (
	"context"
	"os"
	"path/filepath"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

func (fs *KFS) Upload(ctx context.Context, branchName string, dstPath string,
	srcPath string, uploadProcess UploadProcess) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	backupCtx := storage.NewWalkerCtx[sqlite.FileOrDir](ctx, srcPath, &uploadVisitor{
		fs:            fs,
		uploadProcess: uploadProcess,
	})
	scanResp, err := backupCtx.Scan()
	if err != nil {
		return
	}
	info, err := os.Stat(srcPath)
	if err != nil {
		return
	}
	fileOrDir := scanResp.(sqlite.FileOrDir)
	modifyTime := uint64(info.ModTime().UnixNano())
	return fs.Db.UpsertDirItem(ctx, branchName, FormatPath(dstPath), sqlite.DirItem{
		Hash:       fileOrDir.Hash(),
		Name:       filepath.Base(dstPath),
		Mode:       uint64(info.Mode()),
		Size:       fileOrDir.Size(),
		Count:      fileOrDir.Count(),
		TotalCount: fileOrDir.TotalCount(),
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	})
}
