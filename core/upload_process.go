package core

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

type UploadProcess interface {
	New(max int, filename string) UploadProcess
	Close() error
	BeforeContent(hash string, filename string)
	MultiWriter(w io.Writer) io.Writer
}

type EmptyUploadProcess struct {
}

func (process *EmptyUploadProcess) New(max int, filename string) UploadProcess {
	return &EmptyUploadProcess{}
}

func (process *EmptyUploadProcess) Close() error {
	return nil
}

func NewFileByName(process UploadProcess, filename string) (sqlite.File, error) {
	ext := filepath.Ext(filename)
	f, err := os.Open(filename)
	if err != nil {
		return sqlite.File{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return sqlite.File{}, err
	}
	hash := sha256.New()
	w := process.MultiWriter(hash)
	_, err = io.Copy(w, f)
	if err != nil {
		return sqlite.File{}, err
	}
	return sqlite.NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(info.Size()), ext), nil
}

func (process *EmptyUploadProcess) BeforeContent(hash string, filename string) {
}

func (process *EmptyUploadProcess) MultiWriter(w io.Writer) io.Writer {
	return w
}
