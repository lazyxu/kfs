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
	SrcPath    string `json:"srcPath"`
	Concurrent int    `json:"concurrent"`
	Index      int    `json:"index"`
	Label      string `json:"label"`
	FilePath   string `json:"filePath"`
	Size       uint64 `json:"size"`
	StackSize  int    `json:"stackSize"`
	Err        error  `json:"err"`
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
