package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// WriteFileHandle is an open for write handle on a File
type WriteFileHandle struct {
	baseHandle
	node   *File
	closed bool
	offset int64
	opened bool
}

func newWriteFileHandle(node *File) *WriteFileHandle {
	return &WriteFileHandle{
		node: node,
	}
}

func (h *WriteFileHandle) Chmod(mode os.FileMode) error     { return h.node.Chmod(mode) }
func (h *WriteFileHandle) Chown(uid, gid int) error         { return e.ENotImpl }
func (h *WriteFileHandle) Close() error                     { return h.node.Close() }
func (h *WriteFileHandle) Fd() uintptr                      { return 0 }
func (h *WriteFileHandle) Name() string                     { return h.node.Path() }
func (h *WriteFileHandle) Read(b []byte) (n int, err error) { return h.node.Read(b) }
func (h *WriteFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	return h.node.ReadAt(b, h.offset)
}
func (h *WriteFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *WriteFileHandle) Stat() (os.FileInfo, error)                           { return h.node.Stat() }
func (h *WriteFileHandle) Truncate(size int64) error                            { return h.node.Truncate(size) }
func (h *WriteFileHandle) Write(b []byte) (n int, err error)                    { return h.node.Write(b) }
func (h *WriteFileHandle) WriteAt(b []byte, off int64) (n int, err error) {
	return h.node.WriteAt(b, off)
}
func (h *WriteFileHandle) WriteString(s string) (n int, err error) { return h.node.Write([]byte(s)) }
func (h *WriteFileHandle) Node() Node                              { return h.node }
