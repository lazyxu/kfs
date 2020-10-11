package core

import (
	"os"
	"path"

	"github.com/lazyxu/kfs/core/e"
)

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
func (kfs *KFS) MkdirAll(name string, perm os.FileMode) error {
	return e.ENotImpl
}

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
func (kfs *KFS) Open(name string) (*File, error) {
	return kfs.OpenFile(name, os.O_RDONLY, 0)
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
func (kfs *KFS) Create(name string) (*File, error) {
	return kfs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// accessModeMask masks off the read modes from the flags
const accessModeMask = os.O_RDONLY | os.O_WRONLY | os.O_RDWR

// OpenFile a file according to the flags and perm provided
func (kfs *KFS) OpenFile(name string, flags int, perm os.FileMode) (node *File, err error) {
	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if flags&accessModeMask == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		return nil, e.ErrInvalid
	}

	node, err = kfs.GetFile(name)
	if err != nil {
		if err != e.ErrNotExist || flags&os.O_CREATE == 0 {
			return nil, err
		}
		// If not found and O_CREATE then create the file
		dir, leaf, err := kfs.getDirAndLeaf(name)
		if err != nil {
			return nil, err
		}
		node, err = dir.Create(leaf, flags)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (kfs *KFS) Remove(name string) error {
	parent, leaf := path.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return err
	}
	return dir.Remove(leaf)
}
