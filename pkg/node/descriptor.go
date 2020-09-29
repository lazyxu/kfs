package node

import (
	"os"

	"github.com/lazyxu/kfs/kfs/e"
)

type Descriptor interface {
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

	Flush() error
	Release() error
	Node() Node
}

// baseDescriptor implements all the missing methods
type baseDescriptor struct{}

func (h baseDescriptor) Chdir() error                                         { return e.ENotImpl }
func (h baseDescriptor) Chmod(mode os.FileMode) error                         { return e.ENotImpl }
func (h baseDescriptor) Chown(uid, gid int) error                             { return e.ENotImpl }
func (h baseDescriptor) Close() error                                         { return e.ENotImpl }
func (h baseDescriptor) Fd() uintptr                                          { return 0 }
func (h baseDescriptor) Name() string                                         { return "" }
func (h baseDescriptor) Read(b []byte) (n int, err error)                     { return 0, e.ENotImpl }
func (h baseDescriptor) ReadAt(b []byte, off int64) (n int, err error)        { return 0, e.ENotImpl }
func (h baseDescriptor) Readdir(n int) ([]os.FileInfo, error)                 { return nil, e.ENotImpl }
func (h baseDescriptor) Readdirnames(n int) (names []string, err error)       { return nil, e.ENotImpl }
func (h baseDescriptor) Seek(offset int64, whence int) (ret int64, err error) { return 0, e.ENotImpl }
func (h baseDescriptor) Stat() (os.FileInfo, error)                           { return nil, e.ENotImpl }
func (h baseDescriptor) Sync() error                                          { return nil }
func (h baseDescriptor) Truncate(size int64) error                            { return e.ENotImpl }
func (h baseDescriptor) Write(b []byte) (n int, err error)                    { return 0, e.ENotImpl }
func (h baseDescriptor) WriteAt(b []byte, off int64) (n int, err error)       { return 0, e.ENotImpl }
func (h baseDescriptor) WriteString(s string) (n int, err error)              { return 0, e.ENotImpl }
func (h baseDescriptor) Flush() (err error)                                   { return e.ENotImpl }
func (h baseDescriptor) Release() (err error)                                 { return e.ENotImpl }
func (h baseDescriptor) Node() Node                                           { return nil }
