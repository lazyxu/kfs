package core

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/lazyxu/kfs/dao"
	"io"
	"os"
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

func NewFileByName(process UploadProcess, filename string) (dao.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return dao.File{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return dao.File{}, err
	}
	hash := sha256.New()
	w := process.MultiWriter(hash)
	_, err = io.Copy(w, f)
	if err != nil {
		return dao.File{}, err
	}
	return dao.NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(info.Size())), nil
}

func (process *EmptyUploadProcess) BeforeContent(hash string, filename string) {
}

func (process *EmptyUploadProcess) MultiWriter(w io.Writer) io.Writer {
	return w
}
