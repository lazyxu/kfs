package localfs

import "path/filepath"

type IgnoreRule interface {
	Ignore(filename string) bool
}

type AbsPathIgnore struct {
	AbsPath string
}

func (i *AbsPathIgnore) Ignore(filename string) bool {
	return i.AbsPath == filename
}

type FileNameIgnore struct {
	FileName string
}

func (i *FileNameIgnore) Ignore(filename string) bool {
	return i.FileName == filepath.Base(filename)
}
