package core

import (
	"os"
	"time"

	"github.com/lazyxu/kfs/core/e"
)

// DirHandle represents an open directory
type DirHandle struct {
	baseHandle
	node *Dir
}

// newDirHandle opens a directory for read
func newDirHandle(node *Dir) *DirHandle {
	return &DirHandle{
		node: node,
	}
}

func (h DirHandle) Chdir() error                 { return e.ENotImpl }
func (h DirHandle) Chmod(mode os.FileMode) error { return h.node.Chmod(mode) }
func (h DirHandle) Chown(uid, gid int) error     { return e.ENotImpl }
func (h DirHandle) Close() error                 { return h.node.Close() }
func (h DirHandle) Fd() uintptr                  { return 0 }
func (h DirHandle) Name() string                 { return h.node.Name() }
func (h DirHandle) Readdir(n int) ([]os.FileInfo, error) {
	dirs, err := h.node.Readdir(n)
	if err != nil {
		return nil, err
	}
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
func (h DirHandle) Readdirnames(n int) (names []string, err error) { return h.node.Readdirnames(n) }
func (h DirHandle) Stat() (os.FileInfo, error)                     { return h.node.Stat() }
func (h DirHandle) Node() Node                                     { return h.node }
