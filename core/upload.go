package core

import (
	"net"
)

type UploadConfig struct {
	UploadProcess UploadProcess
	Encoder       string
	Concurrent    int
	Verbose       bool
}

type Process struct {
	Index     int
	Label     string
	FilePath  string
	Size      uint64
	StackSize int
	Err       error
}

type UploadProcess interface {
	New(srcPath string, concurrent int, conns []net.Conn) UploadProcess
	Close()
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

func (h *EmptyUploadProcess) Close() {
}

func (h *EmptyUploadProcess) ErrHandler(filePath string, err error) {
	println(filePath+":", err.Error())
}

func (h *EmptyUploadProcess) Verbose() bool {
	return false
}
