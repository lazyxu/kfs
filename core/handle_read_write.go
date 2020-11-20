package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// RWFileHandle is a handle that can be open for read and write.
//
// It will be open to a temporary file which, when closed, will be
// transferred to the remote.
type RWFileHandle struct {
	baseHandle
	node   *File
	closed bool
	offset int64
	opened bool
}

func newRWFileHandle(node *File) *RWFileHandle {
	return &RWFileHandle{
		node: node,
	}
}

func (h *RWFileHandle) Chmod(mode os.FileMode) error     { return h.node.Chmod(mode) }
func (h *RWFileHandle) Chown(uid, gid int) error         { return e.ENotImpl }
func (h *RWFileHandle) Close() error                     { return h.node.Close() }
func (h *RWFileHandle) Fd() uintptr                      { return 0 }
func (h *RWFileHandle) Name() string                     { return h.node.Path() }
func (h *RWFileHandle) Read(b []byte) (n int, err error) { return h.node.Read(b) }
func (h *RWFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	return h.node.ReadAt(b, h.offset)
}
func (h *RWFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *RWFileHandle) Stat() (os.FileInfo, error)                           { return h.node.Stat() }
func (h *RWFileHandle) Truncate(size int64) error                            { return h.node.Truncate(size) }
func (h *RWFileHandle) Write(b []byte) (n int, err error)                    { return h.node.Write(b) }
func (h *RWFileHandle) WriteAt(b []byte, off int64) (n int, err error)       { return h.node.WriteAt(b, off) }
func (h *RWFileHandle) WriteString(s string) (n int, err error)              { return h.node.Write([]byte(s)) }
func (h *RWFileHandle) Node() Node                                           { return h.node }
