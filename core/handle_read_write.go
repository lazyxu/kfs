package core

import (
	"errors"
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
	flags  int
}

func newRWFileHandle(kfs *KFS, path string, flags int) *RWFileHandle {
	return &RWFileHandle{
		kfs:   kfs,
		path:  path,
		flags: flags,
	}
}
func (h *RWFileHandle) appendMode() bool {
	return (h.flags & os.O_APPEND) != 0
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
	h.offset = 0
	h.closed = true
	h.opened = false
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
	n, err = node.ReadAt(b, h.offset)
	if err != nil {
		return 0, wrapErr("read", h.path, err)
	}
	h.offset += int64(n)
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
	return node.ReadAt(b, off)
}
func (h *RWFileHandle) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
		h.offset = offset
	case 1:
		h.offset = h.offset + offset
	case 2:
		node, err := h.Node()
		if err != nil {
			return 0, err
		}
		h.offset = node.Size() + offset
	default:
		return 0, e.ErrInvalid
	}
	return
}

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
	n, err = node.WriteAt(b, h.offset)
	if err != nil {
		return 0, err
	}
	if h.appendMode() {
		h.offset = node.Size()
	}
	h.offset += int64(n)
	return n, nil
}

var ErrWriteAtInAppendMode = errors.New("os: invalid use of WriteAt on file opened with O_APPEND")

func (h *RWFileHandle) WriteAt(b []byte, off int64) (n int, err error) {
	if h.appendMode() {
		return 0, ErrWriteAtInAppendMode
	}
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
	if h.appendMode() {
		h.offset = node.Size()
	}
	n, err = node.WriteAt([]byte(s), h.offset)
	if err != nil {
		return 0, err
	}
	h.offset += int64(n)
	return n, nil
}
func (h *RWFileHandle) Node() (Node, error) { return h.kfs.GetFile(h.path) }
