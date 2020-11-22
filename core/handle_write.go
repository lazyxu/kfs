package core

import (
	"github.com/lazyxu/kfs/core/e"
)

// WriteFileHandle is an open for write handle on a File
type WriteFileHandle struct {
	RWFileHandle
}

func newWriteFileHandle(kfs *KFS, path string) *WriteFileHandle {
	return &WriteFileHandle{RWFileHandle{
		kfs:  kfs,
		path: path,
	}}
}

func (h *WriteFileHandle) Read(b []byte) (n int, err error) {
	return 0, e.ErrPermission
}
func (h *WriteFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, e.ErrPermission
}
