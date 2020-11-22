package core

import (
	"github.com/lazyxu/kfs/core/e"
)

// ReadFileHandle is an open for read file handle on a File
type ReadFileHandle struct {
	RWFileHandle
}

func newReadFileHandle(kfs *KFS, path string) *ReadFileHandle {
	return &ReadFileHandle{RWFileHandle{
		kfs:  kfs,
		path: path,
	}}
}

func (h *ReadFileHandle) Truncate(size int64) error                      { return e.ErrPermission }
func (h *ReadFileHandle) Write(b []byte) (n int, err error)              { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteAt(b []byte, off int64) (n int, err error) { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteString(s string) (n int, err error)        { return 0, e.ErrPermission }
