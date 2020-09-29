package kfs

import (
	"os"
	"path"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/kfs/e"
)

func (kfs *KFS) Remove(filePath string) error {
	parentDir, leaf := path.Split(filePath)
	dir, err := kfs.getNodeDir(parentDir)
	if err != nil {
		return err
	}
	return dir.Remove(leaf)
}

// accessModeMask masks off the read modes from the flags
const accessModeMask = os.O_RDONLY | os.O_WRONLY | os.O_RDWR

// OpenFile a file according to the flags and perm provided
func (kfs *KFS) OpenFile(name string, flags int, perm os.FileMode) (node Item, err error) {
	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if flags&accessModeMask == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		return nil, e.ErrInvalid
	}

	node, err = kfs.getNode(name)
	if err != nil {
		if err != e.ErrNotExist || flags&os.O_CREATE == 0 {
			return nil, err
		}
		// If not found and O_CREATE then create the file
		dir, leaf, err := kfs.getDir(name)
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

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
func (kfs *KFS) Open(name string) (Item, error) {
	return kfs.OpenFile(name, os.O_RDONLY, 0)
}

func (kfs *KFS) Read(path string, buff []byte, off int64) (int64, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return 0, err
	}
	file, ok := n.(*ItemFile)
	if !ok {
		return 0, err
	}
	content, err := file.GetContent()
	if err != nil {
		return 0, err
	}
	end := off + int64(len(buff))
	if end > int64(len(content)) {
		end = int64(len(content))
	}
	if end < off {
		return 0, nil
	}
	size := copy(buff, content[off:end])
	return int64(size), nil
}

func (kfs *KFS) Write(path string, buff []byte, offset int64) (int64, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return 0, err
	}
	file, ok := n.(*ItemFile)
	if !ok {
		return 0, e.ErrNotExist
	}
	return file.SetContent(buff, offset)
}

func (kfs *KFS) Readdir(path string) ([]object.Object, error) {
	n, err := kfs.GetNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*ItemDir)
	if !ok {
		return nil, e.ENotDir
	}
	itemMap, err := dir.ReadDirAll()
	if err != nil {
		return nil, err
	}
	nodes := make([]object.Object, len(itemMap))
	i := 0
	for _, item := range itemMap {
		nodes[i] = item
		i++
	}
	return nodes, nil
}
