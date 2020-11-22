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
	kfs    *KFS
	path   string
	closed bool
	offset int64
	opened bool
}

func newRWFileHandle(kfs *KFS, path string) *RWFileHandle {
	return &RWFileHandle{
		kfs:  kfs,
		path: path,
	}
}

func (h *RWFileHandle) Chmod(mode os.FileMode) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}
func (h *RWFileHandle) Chown(uid, gid int) error { return e.ENotImpl }
func (h *RWFileHandle) Close() error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	h.closed = true
	return node.Close()
}
func (h *RWFileHandle) Fd() uintptr  { return 0 }
func (h *RWFileHandle) Name() string { return h.path }
func (h *RWFileHandle) Read(b []byte) (n int, err error) {
	if h.closed {
		return 0, wrapErr("read", h.path, e.ErrClosed)
	}
	node, err := h.Node()
	if err != nil {
		return 0, wrapErr("read", h.path, err)
	}
	n, err = node.Read(b)
	if err != nil {
		return 0, wrapErr("read", h.path, err)
	}
	return n, nil
}
func (h *RWFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	if h.closed {
		return 0, e.ErrClosed
	}
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.ReadAt(b, h.offset)
}
func (h *RWFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *RWFileHandle) Stat() (os.FileInfo, error) {
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	return node.Stat()
}
func (h *RWFileHandle) Truncate(size int64) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Truncate(size)
}
func (h *RWFileHandle) Write(b []byte) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Write(b)
}
func (h *RWFileHandle) WriteAt(b []byte, off int64) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.WriteAt(b, off)
}
func (h *RWFileHandle) WriteString(s string) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Write([]byte(s))
}
func (h *RWFileHandle) Node() (Node, error) { return h.kfs.GetFile(h.path) }
