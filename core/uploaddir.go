package core

import (
	"fmt"
	"os"
)

type UploadDirConfig struct {
	UploadDirProcess UploadDirProcess
	Encoder          string
	Concurrent       int
	Verbose          bool
}

type UploadDirProcess interface {
	FilePathFilter(filePath string) bool
	OnFileError(filePath string, err error)
	PushFile(info os.FileInfo)
	StartFile(filePath string, info os.FileInfo)
	StartUploadFile(filePath string, info os.FileInfo, hash string)
	EndUploadFile(filePath string, info os.FileInfo)
	SkipFile(filePath string, info os.FileInfo, hash string)
	EndFile(filePath string, info os.FileInfo)
	StartDir(filePath string, info os.FileInfo, n uint64)
	EndDir(filePath string, info os.FileInfo)
}

type DefaultUploadDirProcessor struct{}

func (d *DefaultUploadDirProcessor) FilePathFilter(filePath string) bool {
	fmt.Printf("[Filter] FilePath: %s\n", filePath)
	return false // 默认不过滤
}

func (d *DefaultUploadDirProcessor) OnFileError(filePath string, err error) {
	fmt.Printf("[Error] %s - %v\n", filePath, err)
}

func (d *DefaultUploadDirProcessor) PushFile(info os.FileInfo) {
	fmt.Printf("[Push] File: %s (Size: %d)\n", info.Name(), info.Size())
}

func (d *DefaultUploadDirProcessor) StartFile(filePath string, info os.FileInfo) {
	fmt.Printf("[Start] File: %s\n", filePath)
}

func (d *DefaultUploadDirProcessor) StartUploadFile(filePath string, info os.FileInfo, hash string) {
	fmt.Printf("[Upload Start] %s | Hash: %s\n", filePath, hash)
}

func (d *DefaultUploadDirProcessor) EndUploadFile(filePath string, info os.FileInfo) {
	fmt.Printf("[Upload Complete] %s\n", filePath)
}

func (d *DefaultUploadDirProcessor) SkipFile(filePath string, info os.FileInfo, hash string) {
	fmt.Printf("[Skip] %s | Hash: %s\n", filePath, hash)
}

func (d *DefaultUploadDirProcessor) EndFile(filePath string, info os.FileInfo) {
	fmt.Printf("[End] File: %s\n", filePath)
}

func (d *DefaultUploadDirProcessor) StartDir(filePath string, info os.FileInfo, n uint64) {
	fmt.Printf("[Dir Start] %s | Files: %d\n", filePath, n)
}

func (d *DefaultUploadDirProcessor) EndDir(filePath string, info os.FileInfo) {
	fmt.Printf("[Dir Complete] %s\n", filePath)
}
