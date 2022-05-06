package noncgo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type File struct {
	fileOrDir
	Ext string
}

func NewFile(hash string, size uint64, ext string) File {
	return File{fileOrDir{hash, size}, ext}
}

func NewFileFromBytes(bytes []byte, ext string) File {
	hash := hex.EncodeToString(sha256.New().Sum(bytes))
	return NewFile(hash, uint64(len(bytes)), ext)
}

func (db *SqliteNonCgoDB) WriteFile(ctx context.Context, file File) error {
	_, err := db._db.ExecContext(ctx, `
	INSERT INTO file VALUES (?, ?, ?);
	`, file.hash, file.size, file.Ext)
	return err
}
