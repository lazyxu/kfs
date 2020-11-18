package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// ReadFileHandle is an open for read file handle on a File
type ReadFileHandle struct {
	baseHandle
	node   *File
	closed bool
	offset int64
	opened bool
}

func newReadFileHandle(node *File) *ReadFileHandle {
	return &ReadFileHandle{
		node: node,
	}
}

func (h *ReadFileHandle) Chmod(mode os.FileMode) error     { return h.node.Chmod(mode) }
func (h *ReadFileHandle) Chown(uid, gid int) error         { return e.ENotImpl }
func (h *ReadFileHandle) Close() error                     { return h.node.Close() }
func (h *ReadFileHandle) Fd() uintptr                      { return 0 }
func (h *ReadFileHandle) Name() string                     { return h.node.Name() }
func (h *ReadFileHandle) Read(b []byte) (n int, err error) { return h.node.Read(b) }
func (h *ReadFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	return h.node.ReadAt(b, h.offset)
}
func (h *ReadFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *ReadFileHandle) Stat() (os.FileInfo, error)                           { return h.node.Stat() }
func (h *ReadFileHandle) Truncate(size int64) error                            { return e.ErrPermission }
func (h *ReadFileHandle) Write(b []byte) (n int, err error)                    { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteAt(b []byte, off int64) (n int, err error)       { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteString(s string) (n int, err error)              { return 0, e.ErrPermission }
func (h *ReadFileHandle) Node() Node                                           { return h.node }
