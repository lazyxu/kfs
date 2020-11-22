package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// ReadFileHandle is an open for read file handle on a File
type ReadFileHandle struct {
	baseHandle
	kfs    *KFS
	path   string
	closed bool
	offset int64
	opened bool
}

func newReadFileHandle(kfs *KFS, path string) *ReadFileHandle {
	return &ReadFileHandle{
		kfs:  kfs,
		path: path,
	}
}

func (h *ReadFileHandle) Chmod(mode os.FileMode) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}
func (h *ReadFileHandle) Chown(uid, gid int) error { return e.ENotImpl }
func (h *ReadFileHandle) Close() error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Close()
}
func (h *ReadFileHandle) Fd() uintptr  { return 0 }
func (h *ReadFileHandle) Name() string { return h.path }
func (h *ReadFileHandle) Read(b []byte) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.Read(b)
}
func (h *ReadFileHandle) ReadAt(b []byte, off int64) (n int, err error) {
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.ReadAt(b, h.offset)
}
func (h *ReadFileHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h *ReadFileHandle) Stat() (os.FileInfo, error) {
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	return node.Stat()
}
func (h *ReadFileHandle) Truncate(size int64) error                      { return e.ErrPermission }
func (h *ReadFileHandle) Write(b []byte) (n int, err error)              { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteAt(b []byte, off int64) (n int, err error) { return 0, e.ErrPermission }
func (h *ReadFileHandle) WriteString(s string) (n int, err error)        { return 0, e.ErrPermission }
func (h *ReadFileHandle) Node() (Node, error)                            { return h.kfs.GetFile(h.path) }
