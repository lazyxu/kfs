package core

import (
	"os"
	"path"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
)

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
func (kfs *KFS) Mkdir(name string, perm os.FileMode) error {
	parent, leaf := path.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return err
	}
	return dir.add(object.NewDirMetadata(leaf, perm), object.EmptyDir)
}

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
func (kfs *KFS) Open(name string) (Handle, error) {
	return kfs.OpenFile(name, os.O_RDONLY, 0)
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
func (kfs *KFS) Create(name string) (Handle, error) {
	return kfs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// accessModeMask masks off the read modes from the flags
const accessModeMask = os.O_RDONLY | os.O_WRONLY | os.O_RDWR

// OpenFile a file according to the flags and perm provided
func (kfs *KFS) OpenFile(name string, flags int, perm os.FileMode) (h Handle, err error) {
	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if flags&accessModeMask == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		return nil, e.ErrInvalid
	}
	var node Node
	node, err = kfs.getNode(name)
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
	return node.Open(flags)
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
func (kfs *KFS) Rename(oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := kfs.GetDir(oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.get(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := kfs.GetDir(newParent)
	if err != nil {
		return err
	}
	newMetadata, err := newDir.get(newName)
	if err == e.ErrNotExist {
		err := oldDir.remove(oldName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return kfs.move(&metadata, newDir)
	}
	if err != nil {
		return err
	}
	if oldMetadata.IsFile() && newMetadata.IsFile() {
		err = oldDir.remove(oldName, true)
		if err != nil {
			return err
		}
		err = newDir.remove(newName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return kfs.move(&metadata, newDir)
	}
	if newMetadata.IsDir() {
		return e.EIsDir
	}
	return nil
}

func (kfs *KFS) move(metadata *object.Metadata, newDir *Dir) error {
	if metadata.IsFile() {
		blob := new(object.Blob)
		err := blob.Read(kfs.scheduler, metadata.Hash)
		if err != nil {
			return err
		}
		return newDir.add(metadata, blob)
	} else {
		tree := new(object.Tree)
		err := tree.Read(kfs.scheduler, metadata.Hash)
		if err != nil {
			return err
		}
		return newDir.add(metadata, tree)
	}
}

// Remove removes the named file or (empty) directory.
func (kfs *KFS) Remove(name string) error {
	parent, leaf := path.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return err
	}
	return dir.remove(leaf, false)
}

// Chmod changes the mode of the named file to mode.
func (kfs *KFS) Chmod(name string, mode os.FileMode) error {
	node, err := kfs.getNode(name)
	if err != nil {
		return err
	}
	return node.Chmod(mode)
}

// Chdir changes the current working directory to the named directory.
func (kfs *KFS) Chdir(dir string) error {
	node, err := kfs.GetDir(dir)
	if err != nil {
		return err
	}
	kfs.pwd = node.Path()
	return nil
}
