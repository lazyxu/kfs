package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/object"
)

var kfs = New(&kfscommon.Options{
	UID:       uint32(os.Getuid()),
	GID:       uint32(os.Getgid()),
	DirPerms:  object.S_IFDIR | 0755,
	FilePerms: object.S_IFREG | 0644,
})

func Open(name string) (Node, error) {
	return kfs.Open(name)
}

func Create(name string) (Node, error) {
	return kfs.Create(name)
}

func Stat(name string) (os.FileInfo, error) {
	return kfs.Stat(name)
}

func Lstat(name string) (os.FileInfo, error) {
	return kfs.Lstat(name)
}

// Symlink creates newname as a symbolic link to oldname.
func Symlink(oldname, newname string) error {
	return e.ENotImpl
}
