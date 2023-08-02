package core

import (
	"github.com/lazyxu/kfs/dao"
	"net"
	"os"
)

type UploadConfig struct {
	UploadProcess UploadProcess
	Encoder       string
	Concurrent    int
	Verbose       bool
}

type Process struct {
	SrcPath    string
	Concurrent int
	Index      int
	Label      string
	FilePath   string
	Size       uint64
	StackSize  int
	Err        error
}

type FileResp struct {
	FileOrDir dao.IFileOrDir
	Info      os.FileInfo
}

type UploadProcess interface {
	New(srcPath string, concurrent int, conns []net.Conn) UploadProcess
	Close(resp FileResp, err error)
	StackSizeHandler(size int)
	ErrHandler(filePath string, err error)
	Show(p *Process)
	Verbose() bool
}

type EmptyUploadProcess struct {
}

func (h *EmptyUploadProcess) Show(p *Process) {
}

func (h *EmptyUploadProcess) StackSizeHandler(size int) {
}

func (h *EmptyUploadProcess) New(srcPath string, concurrent int, conns []net.Conn) UploadProcess {
	return h
}

func (h *EmptyUploadProcess) Close(resp FileResp, err error) {
}

func (h *EmptyUploadProcess) ErrHandler(filePath string, err error) {
	println(filePath+":", err.Error())
}

func (h *EmptyUploadProcess) Verbose() bool {
	return false
}
