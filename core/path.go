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

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (kfs *KFS) RemoveAll(name string) error {
	parent, leaf := path.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return err
	}
	return dir.remove(leaf, true)
}
