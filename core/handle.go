package core

import (
	"os"
	"time"

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

// Handle is the interface satisfied by open files or directories.
// It is the methods on *os.File, plus a few more useful for FUSE
// filing systems.  Not all of them are supported.
type Handle interface {
	OsFiler
	// Additional methods useful for FUSE filesystems
	Flush() error
	Release() error
	Node() (Node, error)
}

// baseHandle implements all the missing methods
type baseHandle struct{}

func (h baseHandle) Chdir() error                                         { return e.ENotImpl }
func (h baseHandle) Chmod(mode os.FileMode) error                         { return e.ENotImpl }
func (h baseHandle) Chown(uid, gid int) error                             { return e.ENotImpl }
func (h baseHandle) Close() error                                         { return e.ENotImpl }
func (h baseHandle) Fd() uintptr                                          { return 0 }
func (h baseHandle) Name() string                                         { return "" }
func (h baseHandle) Read(b []byte) (n int, err error)                     { return 0, e.ENotImpl }
func (h baseHandle) ReadAt(b []byte, off int64) (n int, err error)        { return 0, e.ENotImpl }
func (h baseHandle) Readdir(n int) ([]os.FileInfo, error)                 { return nil, e.ENotImpl }
func (h baseHandle) Readdirnames(n int) (names []string, err error)       { return nil, e.ENotImpl }
func (h baseHandle) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h baseHandle) Stat() (os.FileInfo, error)                           { return nil, e.ENotImpl }
func (h baseHandle) Sync() error                                          { return nil }
func (h baseHandle) Truncate(size int64) error                            { return e.ENotImpl }
func (h baseHandle) Write(b []byte) (n int, err error)                    { return 0, e.ENotImpl }
func (h baseHandle) WriteAt(b []byte, off int64) (n int, err error)       { return 0, e.ENotImpl }
func (h baseHandle) WriteString(s string) (n int, err error)              { return 0, e.ENotImpl }
func (h baseHandle) Flush() (err error)                                   { return e.ENotImpl }
func (h baseHandle) Release() (err error)                                 { return e.ENotImpl }
func (h baseHandle) Node() (Node, error)                                  { return nil, e.ENotImpl }

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
