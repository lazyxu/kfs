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
	FilePathFilter(filePath string) bool
	OnFileError(filePath string, err error)
	PushFile(info os.FileInfo)
	StartFile(filePath string, info os.FileInfo)
	StartUploadFile(filePath string, info os.FileInfo, hash string)
	EndUploadFile(filePath string, info os.FileInfo)
	SkipFile(filePath string, info os.FileInfo, hash string)
	EndFile(filePath string, info os.FileInfo)
	StartDir(filePath string, n uint64)
	EndDir(filePath string, info os.FileInfo)
}
