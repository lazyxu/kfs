package local

import (
	"context"
	"fmt"
	"os"
	"time"

	storage "github.com/lazyxu/kfs/storage/local"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	b *Backup
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
		_, err = v.b.s.Write(file.Hash(), f)
		if err != nil {
			return nil, err
		}
		err = v.b.db.WriteFile(ctx, file)
		fmt.Printf("upload file %s %+v\n", filename, file)
		return file, err
	} else if info.IsDir() {
		dirItems := make([]sqlite.DirItem, len(infos))
		for i, info := range infos {
			name := info.Name()
			now := uint64(time.Now().UnixNano())
			modifyTime := uint64(info.ModTime().UnixNano())
			dirItems[i] = sqlite.NewDirItem(rets[i], name, uint64(info.Mode()), now, modifyTime, now, now)
		}
		dir, err := v.b.db.WriteDir(ctx, dirItems)
		if err != nil {
			return nil, err
		}
		fmt.Printf("upload dir %s %+v\n", filename, dir)
		return dir, nil
	}
	return nil, nil
}
