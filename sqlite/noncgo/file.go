package noncgo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
)

type File struct {
	fileOrDir
	Ext string
}

func NewFile(hash string, size uint64, ext string) File {
	return File{fileOrDir{hash, size}, ext}
}

func NewFileByBytes(bytes []byte, ext string) File {
	hash := sha256.New()
	hash.Write(bytes)
	return NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(len(bytes)), ext)
}

func NewFileByName(filename string) (File, error) {
	ext := filepath.Ext(filename)
	f, err := os.Open(filename)
	if err != nil {
		return File{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return File{}, err
	}
	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return File{}, err
	}
	return NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(info.Size()), ext), nil
}

func (db *DB) WriteFile(ctx context.Context, file File) error {
	return db.writeFile(ctx, db._db, file)
}

func (db *DB) writeFile(ctx context.Context, txOrDb TxOrDb, file File) error {
	// TODO: update ext if duplicated
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO file VALUES (?, ?, ?);
	`, file.hash, file.size, file.Ext)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}

func (db *DB) UploadFile(ctx context.Context, branchName string, splitPath []string, hash string, size uint64,
	mode uint64, createTime uint64, modifyTime uint64, changeTime uint64, accessTime uint64) (commit Commit, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err1 := tx.Rollback()
				if err1 != nil {
					panic(err1) // should not happen
				}
				return
			}
		}
	}()
	return db.updateDirItem(ctx, tx, branchName, splitPath, func(dirItemsList [][]DirItem) error {
		i := len(dirItemsList) - 1
		find := false
		name := splitPath[len(splitPath)-1]
		ext := path.Ext(name)
		file := NewFile(hash, size, ext)
		err = db.writeFile(ctx, tx, file)
		if err != nil {
			return err
		}
		newItem := NewDirItem(file, name, mode, createTime, modifyTime, changeTime, accessTime)
		for j, dirItem := range dirItemsList[i] {
			if dirItem.Name == splitPath[i] {
				dirItemsList[i][j] = newItem
				find = true
				break
			}
		}
		if !find {
			dirItemsList[i] = append(dirItemsList[i], newItem)
		}
		return nil
	})
}
