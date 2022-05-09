package local

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lazyxu/kfs/storage/local"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type uploadVisitor struct {
	local.EmptyVisitor
	backup *Backup
}

func (v *uploadVisitor) HasExit() bool {
	return true
}

func (v *uploadVisitor) Exit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []any) (any, error) {
	if info.Mode().IsRegular() {
		file, err := sqlite.NewFileByName(filename)
		if err != nil {
			return nil, err
		}
		err = v.backup.db.WriteFile(ctx, file)
		fmt.Printf("upload file %s %+v\n", filename, file)
		return file, err
	} else if info.IsDir() {
		dirItems := make([]sqlite.DirItem, len(infos))
		for i, info := range infos {
			name := info.Name()
			now := uint64(time.Now().UnixNano())
			modifyTime := uint64(info.ModTime().UnixNano())
			if file, ok := rets[i].(sqlite.File); ok && info.Mode().IsRegular() {
				dirItems[i] = sqlite.NewDirItem(file, name, uint64(info.Mode()), now, modifyTime, now, now, "")
			} else if dir, ok := rets[i].(sqlite.Dir); ok && info.IsDir() {
				dirItems[i] = sqlite.NewDirItem(dir, name, uint64(info.Mode()), now, modifyTime, now, now, "")
			}
		}
		dir, err := v.backup.db.WriteDir(ctx, dirItems)
		if err != nil {
			return nil, err
		}
		fmt.Printf("upload dir %s %+v\n", filename, dir)
		return dir, nil
	}
	return nil, nil
}
