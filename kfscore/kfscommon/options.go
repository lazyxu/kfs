package kfscommon

import (
	"os"
)

// Options is options for creating the vfs
type Options struct {
	UID       uint32
	GID       uint32
	DirPerms  os.FileMode
	FilePerms os.FileMode
}
