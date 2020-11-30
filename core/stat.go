package core

import (
	"os"
)

// Stat returns a FileInfo describing the named file.
func (kfs *KFS) Stat(name string) (os.FileInfo, error) {
	n, err := kfs.GetNode(name)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
func (kfs *KFS) Lstat(name string) (os.FileInfo, error) {
	return kfs.Stat(name)
}
