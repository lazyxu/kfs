package core

import (
	"context"
	"os"

	storage "github.com/lazyxu/kfs/storage/local"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	fs *KFS
}

func (v *uploadVisitor) HasExit() bool {
	return true
}

func (v *uploadVisitor) Exit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []sqlite.FileOrDir) (sqlite.FileOrDir, error) {
	if info.Mode().IsRegular() {
		file, err := sqlite.NewFileByName(filename)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = v.fs.S.Write(file.Hash(), f)
		if err != nil {
			return nil, err
		}
		err = v.fs.db.WriteFile(ctx, file)
		return file, err
	} else if info.IsDir() {
		dirItems := make([]sqlite.DirItem, len(infos))
		for i, info := range infos {
			if rets[i] == nil {
				continue
			}
			name := info.Name()
			modifyTime := uint64(info.ModTime().UnixNano())
			dirItems[i] = sqlite.NewDirItem(rets[i], name, uint64(info.Mode()), modifyTime, modifyTime, modifyTime, modifyTime)
		}
		dir, err := v.fs.db.WriteDir(ctx, dirItems)
		if err != nil {
			return nil, err
		}
		return dir, nil
	}
	return nil, nil
}
