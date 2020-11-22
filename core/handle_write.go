package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// WriteFileHandle is an open for write handle on a File
type WriteFileHandle struct {
	baseHandle
	kfs    *KFS
	path   string
	closed bool
	offset int64
	opened bool
}

func newWriteFileHandle(kfs *KFS, path string) *WriteFileHandle {
	return &WriteFileHandle{
		kfs:  kfs,
		path: path,
	}
}

func (h *WriteFileHandle) Chmod(mode os.FileMode) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}
func (h *WriteFileHandle) Chown(uid, gid int) error { return e.ENotImpl }
func (h *WriteFileHandle) Close() error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Close()
}
func (h *WriteFileHandle) Fd() uintptr  { return 0 }
func (h *WriteFileHandle) Name() string { return h.path }
func (h *WriteFileHandle) Read(b []byte) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Read(b)
}
func (h *WriteFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.ReadAt(b, h.offset)
}
func (h *WriteFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *WriteFileHandle) Stat() (os.FileInfo, error) {
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	return node.Stat()
}
func (h *WriteFileHandle) Truncate(size int64) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Truncate(size)
}
func (h *WriteFileHandle) Write(b []byte) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Write(b)
}
func (h *WriteFileHandle) WriteAt(b []byte, off int64) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.WriteAt(b, off)
}
func (h *WriteFileHandle) WriteString(s string) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Write([]byte(s))
}
func (h *WriteFileHandle) Node() (Node, error) { return h.kfs.GetFile(h.path) }
