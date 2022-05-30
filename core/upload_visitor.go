package core

import (
	"context"
	"io"
	"os"
	"path/filepath"

	storage "github.com/lazyxu/kfs/storage/local"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	fs            *KFS
	uploadProcess UploadProcess
}

func (v *uploadVisitor) HasExit() bool {
	return true
}

func (v *uploadVisitor) Exit(ctx context.Context, filePath string, info os.FileInfo, infos []os.FileInfo, rets []sqlite.FileOrDir) (sqlite.FileOrDir, error) {
	if info.Mode().IsRegular() {
		v.uploadProcess = v.uploadProcess.New(int(info.Size()), filepath.Base(filePath))
		defer v.uploadProcess.Close()
		file, err := NewFileByName(v.uploadProcess, filePath)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = v.fs.S.WriteFn(file.Hash(), func(w io.Writer, hasher io.Writer) error {
			rr := io.TeeReader(f, hasher)
			v.uploadProcess.BeforeContent(file.Hash(), filepath.Base(filePath))
			w = v.uploadProcess.MultiWriter(w)
			_, err := io.Copy(w, rr)
			return err
		})
		if err != nil {
			return nil, err
		}
		err = v.fs.Db.WriteFile(ctx, file)
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
		dir, err := v.fs.Db.WriteDir(ctx, dirItems)
		if err != nil {
			return nil, err
		}
		return dir, nil
	}
	return nil, nil
}
