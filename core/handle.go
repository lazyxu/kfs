package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lazyxu/kfs/node"

	"github.com/lazyxu/kfs/core/e"
)

// Check interfaces
var (
	_ OsFiler = (*os.File)(nil)
)

// OsFiler is the methods on *os.File
type OsFiler interface {
	Chdir() error
	Chmod(mode os.FileMode) error
	Chown(uid, gid int) error
	Close() error
	Fd() uintptr
	Name() string
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Readdir(n int) ([]os.FileInfo, error)
	Readdirnames(n int) (names []string, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteString(s string) (n int, err error)
}

// handle is the interface satisfied by open files or directories.
// It is the methods on *os.File, plus a few more useful for FUSE
// filing systems.  Not all of them are supported.
type handle interface {
	// Additional methods useful for FUSE filesystems
	Flush() error
	Release() error
	Node() (node.Node, error)
}

// baseHandle implements all the missing methods
type Handle struct {
	kfs        *KFS
	path       string
	offset     int64
	nameOffset int
	isDir      bool
	closed     bool
	opened     bool
	read       bool
	write      bool
	append     bool
}

func (h *Handle) Chdir() error {
	if h == nil {
		return e.ErrInvalid
	}
	h.kfs.pwd = h.path
	return nil
}

func (h *Handle) Chmod(mode os.FileMode) error {
	if h == nil {
		return e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}

func (h *Handle) Chown(uid, gid int) error { return e.ErrInvalid }
func (h *Handle) Close() error {
	if h == nil {
		return e.ErrInvalid
	}
	if h.closed {
		return wrapErr("close", h.path, e.ErrClosed)
	}
	node, err := h.Node()
	if err != nil {
		return err
	}
	h.offset = 0
	h.nameOffset = 0
	h.closed = true
	h.opened = false
	return node.Close()
}
func (h *Handle) Fd() uintptr  { return 0 }
func (h *Handle) Name() string { return h.path }
func (h *Handle) Read(b []byte) (n int, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
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
func (h *Handle) ReadAt(b []byte, off int64) (n int, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
	if h.closed {
		return 0, e.ErrClosed
	}
	if off < 0 {
		return 0, fmt.Errorf("negative offset: %d", off)
	}
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.ReadAt(b, off)
}
func (h *Handle) Readdir(n int) ([]os.FileInfo, error) {
	if h == nil {
		return nil, e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	dirs, err := node.Readdir(n, int(h.offset))
	if err != nil {
		return []os.FileInfo{}, err
	}
	h.offset += int64(len(dirs))
	infos := make([]os.FileInfo, len(dirs))
	for i, dir := range dirs {
		infos[i] = &fileInfo{
			name:    dir.Name(),
			size:    dir.Size(),
			mode:    dir.Mode(),
			modTime: dir.ModifyTime(),
		}
	}
	return infos, nil
}
func (h *Handle) Readdirnames(n int) ([]string, error) {
	if h == nil {
		return nil, e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	names, err := node.Readdirnames(n, h.nameOffset)
	if err != nil {
		return names, err
	}
	h.nameOffset += len(names)
	return names, err
}
func (h *Handle) Seek(offset int64, whence int) (ret int64, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
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

func (h *Handle) Stat() (os.FileInfo, error) {
	if h == nil {
		return nil, e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	return node.Stat()
}
func (h *Handle) Sync() error {
	return e.ErrInvalid
}
func (h *Handle) Truncate(size int64) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Truncate(size)
}
func (h *Handle) Write(b []byte) (n int, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	n, err = node.WriteAt(b, h.offset)
	if err != nil {
		return 0, err
	}
	if h.append {
		h.offset = node.Size()
	}
	h.offset += int64(n)
	return n, nil
}

var ErrWriteAtInAppendMode = errors.New("os: invalid use of WriteAt on file opened with O_APPEND")

func (h *Handle) WriteAt(b []byte, off int64) (n int, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
	if h.append {
		return 0, ErrWriteAtInAppendMode
	}
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	return node.WriteAt(b, off)
}
func (h *Handle) WriteString(s string) (n int, err error) {
	if h == nil {
		return 0, e.ErrInvalid
	}
	node, err := h.Node()
	if err != nil {
		return 0, err
	}
	if h.append {
		h.offset = node.Size()
	}
	n, err = node.WriteAt([]byte(s), h.offset)
	if err != nil {
		return 0, err
	}
	h.offset += int64(n)
	return n, nil
}

func (h *Handle) Node() (node.Node, error) {
	if h == nil {
		return nil, e.ErrInvalid
	}
	return h.kfs.GetNode(h.path)
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return f.size }
func (f *fileInfo) Mode() os.FileMode  { return f.mode }
func (f *fileInfo) ModTime() time.Time { return f.modTime }
func (f *fileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f *fileInfo) Sys() interface{}   { return nil }

// wrapErr wraps an error that occurred during an operation on an open file.
// It passes io.EOF through unchanged, otherwise converts
// poll.ErrFileClosing to ErrClosed and wraps the error in a PathError.
func wrapErr(op string, path string, err error) error {
	if err == nil || err == io.EOF {
		return err
	}
	return &PathError{op, path, err}
}
