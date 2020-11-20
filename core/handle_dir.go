package core

import (
	"os"
	"time"

	"github.com/lazyxu/kfs/core/e"
)

// DirHandle represents an open directory
type DirHandle struct {
	baseHandle
	node       *Dir
	offset     int
	nameOffset int
}

// newDirHandle opens a directory for read
func newDirHandle(node *Dir) *DirHandle {
	return &DirHandle{
		node: node,
	}
}

func (h *DirHandle) Chdir() error                 { return e.ENotImpl }
func (h *DirHandle) Chmod(mode os.FileMode) error { return h.node.Chmod(mode) }
func (h *DirHandle) Chown(uid, gid int) error     { return e.ENotImpl }
func (h *DirHandle) Close() error {
	h.offset = 0
	h.nameOffset = 0
	return h.node.Close()
}
func (h *DirHandle) Fd() uintptr  { return 0 }
func (h *DirHandle) Name() string { return h.node.Path() }
func (h *DirHandle) Readdir(n int) ([]os.FileInfo, error) {
	dirs, err := h.node.Readdir(n, h.offset)
	if err != nil {
		return []os.FileInfo{}, err
	}
	h.offset += len(dirs)
	infos := make([]os.FileInfo, len(dirs))
	for i, dir := range dirs {
		infos[i] = &fileInfo{
			name:    dir.Name,
			size:    dir.Size,
			mode:    dir.Mode,
			modTime: time.Unix(0, dir.ModifyTime),
		}
	}
	return infos, nil
}
func (h *DirHandle) Readdirnames(n int) ([]string, error) {
	names, err := h.node.Readdirnames(n, h.nameOffset)
	if err != nil {
		return names, err
	}
	h.nameOffset += len(names)
	return names, err
}
func (h *DirHandle) Stat() (os.FileInfo, error) { return h.node.Stat() }
func (h *DirHandle) Node() Node                 { return h.node }
