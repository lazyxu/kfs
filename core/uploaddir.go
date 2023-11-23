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
	StartFile(filePath string, info os.FileInfo)
	EndFile(filePath string, info os.FileInfo)
	StartDir(filePath string, n uint64)
	EndDir(filePath string, info os.FileInfo)
	OnFileError(filePath string, err error)
	PushFile(info os.FileInfo)
}
