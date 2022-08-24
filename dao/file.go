package dao

import (
	"crypto/sha256"
	"encoding/hex"
)

type File struct {
	FileOrDir
}

func NewFile(hash string, size uint64) File {
	return File{FileOrDir{hash, size}}
}

func NewFileByBytes(bytes []byte) File {
	hash := sha256.New()
	hash.Write(bytes)
	return NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(len(bytes)))
}
