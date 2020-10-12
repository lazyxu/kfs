package core

import (
	"os"

	"github.com/lazyxu/kfs/core/e"
)

// Stat returns a FileInfo describing the named file.
func (kfs *KFS) Stat(name string) (os.FileInfo, error) {
	n, err := kfs.getNode(name)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}
	return n, nil
}

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
func (kfs *KFS) Lstat(name string) (os.FileInfo, error) {
	// There is no symbolic link in koala file system.
	return nil, e.ENotImpl
}
