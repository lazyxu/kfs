package core

import (
	"os"
	"time"

	"github.com/lazyxu/kfs/core/e"
)

// DirHandle represents an open directory
type DirHandle struct {
	baseHandle
	kfs        *KFS
	path       string
	offset     int
	nameOffset int
}

// newDirHandle opens a directory for read
func newDirHandle(kfs *KFS, path string) *DirHandle {
	return &DirHandle{
		kfs:  kfs,
		path: path,
	}
}

func (h *DirHandle) Chdir() error { h.kfs.pwd = h.path; return nil }
func (h *DirHandle) Chmod(mode os.FileMode) error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}
func (h *DirHandle) Chown(uid, gid int) error { return e.ENotImpl }
func (h *DirHandle) Close() error {
	node, err := h.Node()
	if err != nil {
		return err
	}
	h.offset = 0
	h.nameOffset = 0
	return node.Close()
}
func (h *DirHandle) Fd() uintptr  { return 0 }
func (h *DirHandle) Name() string { return h.path }
func (h *DirHandle) Readdir(n int) ([]os.FileInfo, error) {
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	dirs, err := node.Readdir(n, h.offset)
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

func (h *DirHandle) Stat() (os.FileInfo, error) {
	node, err := h.Node()
	if err != nil {
		return nil, err
	}
	return node.Stat()
}

func (h *DirHandle) Node() (Node, error) { return h.kfs.GetDir(h.path) }
