package core

import (
	"os"
)

type UploadDirConfig struct {
	UploadDirProcess UploadDirProcess
	Encoder          string
	Concurrent       int
	Verbose          bool
}

type UploadDirProcess interface {
	Close(resp FileResp, err error)
	StackSizeHandler(size int)
	Show(p *Process)
	Verbose() bool
	StartFile(filePath string, info os.FileInfo)
	OnFileError(filePath string, err error)
	EndFile(filePath string, info os.FileInfo)
	PushFile(info os.FileInfo)
	HasPushedAllToStack()
}

type EmptyUploadDirProcess struct {
}

func (h *EmptyUploadDirProcess) Show(p *Process) {
}

func (h *EmptyUploadDirProcess) StackSizeHandler(size int) {
}

func (h *EmptyUploadDirProcess) Close(resp FileResp, err error) {
}

func (h *EmptyUploadDirProcess) Verbose() bool {
	return false
}

func (h *EmptyUploadDirProcess) StartFile(filePath string, info os.FileInfo) {
}

func (h *EmptyUploadDirProcess) OnFileError(filePath string, err error) {
	println(filePath+":", err.Error())
}

func (h *EmptyUploadDirProcess) EndFile(filePath string, info os.FileInfo, exist bool) {
}

func (h *EmptyUploadDirProcess) PushFile(info os.FileInfo) {
}

func (h *EmptyUploadDirProcess) HasPushedAllToStack() {
}
