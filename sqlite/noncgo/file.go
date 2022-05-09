package noncgo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
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

func (db *SqliteNonCgoDB) WriteFile(ctx context.Context, file File) error {
	// TODO: update ext if duplicated
	_, err := db._db.ExecContext(ctx, `
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
