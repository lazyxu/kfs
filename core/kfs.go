package core

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
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
func (kfs *KFS) OpenFile(name string, flags int, perm os.FileMode) (node Node, err error) {
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
func (kfs *KFS) Open(name string) (Node, error) {
	return kfs.OpenFile(name, os.O_RDONLY, 0)
}

func (kfs *KFS) Read(path string, buff []byte, off int64) (int64, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return 0, err
	}
	file, ok := n.(*File)
	if !ok {
		return 0, err
	}
	reader, err := file.GetContent()
	if err != nil {
		return 0, err
	}
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		if err != nil {
			return n, err
		}
	default:
		n, err := io.CopyN(ioutil.Discard, r, off)
		if err != nil {
			return n, err
		}
	}
	num, err := reader.Read(buff)
	return int64(num), err
}

func (kfs *KFS) Write(path string, buff []byte, offset int64) (int64, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return 0, err
	}
	file, ok := n.(*File)
	if !ok {
		return 0, e.ErrNotExist
	}
	return file.SetContent(buff, offset)
}

func (kfs *KFS) Readdir(path string) ([]object.Metadata, error) {
	n, err := kfs.GetNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*Dir)
	if !ok {
		return nil, e.ENotDir
	}
	itemMap, err := dir.ReadDirAll()
	if err != nil {
		return nil, err
	}
	nodes := make([]object.Metadata, len(itemMap))
	i := 0
	for _, item := range itemMap {
		nodes[i] = item
		i++
	}
	return nodes, nil
}
