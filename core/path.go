package core

import (
	"os"
	"path"
	"strings"

	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/object"
)

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
func (kfs *KFS) MkdirAll(path string, perm os.FileMode) error {
	path = strings.Trim(path, "/")
	var node Node
	node = kfs.root
	for path != "" {
		i := strings.IndexRune(path, '/')
		var name string
		if i < 0 {
			name, path = path, ""
		} else {
			name, path = path[:i], path[i+1:]
		}
		if name == "" {
			continue
		}
		dir, ok := node.(*Dir)
		if !ok {
			// We need to look in a directory, but found a file
			return e.ENotDir
		}
		node, ok = dir.items[name]
		if ok {
			continue
		}

		d, err := object.ReadDir(kfs.scheduler, dir.Metadata.Hash)
		if err != nil {
			return err
		}
		metadata, err := d.GetNode(name)
		if err == e.ENoSuchFileOrDir {
			node = &Dir{
				ItemBase: ItemBase{
					kfs:      kfs,
					parent:   dir,
					Metadata: object.NewDirMetadata(name, perm),
				},
				items: make(map[string]Node),
			}
			dir.items[name] = node
			continue
		}
		if err != nil {
			return err
		}
		if metadata.IsDir() {
			node = &Dir{
				ItemBase: ItemBase{
					kfs:      kfs,
					parent:   dir,
					Metadata: metadata,
				},
				items: make(map[string]Node),
			}
			dir.items[name] = node
		} else {
			node = &File{
				ItemBase: ItemBase{
					kfs:      kfs,
					parent:   dir,
					Metadata: metadata,
				},
			}
			dir.items[name] = node
		}
	}
	return nil
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
