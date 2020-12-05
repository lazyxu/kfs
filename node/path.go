package node

import (
	"math"
	"path"
	"strconv"
	"strings"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
)

// GetNode finds the Node by path.
func GetNode(n Node, path string) (Node, error) {
	obj := n.Obj()
	storage := n.Storage()
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
		dir, ok := n.(*Dir)
		if !ok {
			// We need to look in a directory, but found a file
			return nil, e.ENotDir
		}
		n, ok = dir.Items[name]
		if ok {
			continue
		}

		d, err := obj.ReadDir(storage, dir.Metadata.Hash)
		if err != nil {
			return nil, err
		}
		metadata, err := d.GetNode(name)
		if err != nil {
			return nil, err
		}
		if metadata.IsDir() {
			n = NewDir(storage, obj, metadata, dir)
			dir.Items[name] = n
		} else {
			n = NewFile(storage, obj, metadata, dir)
			dir.Items[name] = n
		}
	}
	return n, nil
}

func GetDir(n Node, path string) (dir *Dir, err error) {
	n, err = GetNode(n, path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir, nil
}

func GetFile(n Node, path string) (*File, error) {
	n, err := GetNode(n, path)
	if err != nil {
		return nil, err
	}
	file, ok := n.(*File)
	if !ok {
		return nil, e.ENotFile
	}
	return file, nil
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
func Rename(n Node, oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := GetDir(n, oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := GetDir(n, newParent)
	if err != nil {
		return err
	}
	newMetadata, err := newDir.GetChild(newName)
	if err == e.ErrNotExist {
		err := oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return move(newDir, &metadata)
	}
	if err != nil {
		return err
	}
	if oldMetadata.IsFile() && newMetadata.IsFile() {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		err = newDir.RemoveChild(newName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return move(newDir, &metadata)
	}
	if newMetadata.IsDir() {
		return e.EIsDir
	}
	return nil
}

func Mv(n Node, oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := GetDir(n, oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := GetDir(n, newParent)
	if err != nil {
		return err
	}
	newMetadata, err := newDir.GetChild(newName)
	if err == e.ErrNotExist {
		err := oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return move(newDir, &metadata)
	}
	if err != nil {
		return err
	}
	if oldMetadata.IsFile() && newMetadata.IsFile() {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		err = newDir.RemoveChild(newName, true)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = newName
		return move(newDir, &metadata)
	}
	if newMetadata.IsDir() {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		newDir, err := GetDir(n, newPath)
		if err != nil {
			return err
		}
		metadata := *oldMetadata
		metadata.Name = oldName
		return move(newDir, &metadata)
	}
	return nil
}

func move(newDir *Dir, metadata *object.Metadata) error {
	if metadata.IsFile() {
		blob := newDir.Obj().NewBlob()
		err := blob.Read(newDir.Storage(), metadata.Hash)
		if err != nil {
			return err
		}
		err = newDir.AddChild(metadata, blob)
		return err
	} else {
		tree := newDir.Obj().NewTree()
		err := tree.Read(newDir.Storage(), metadata.Hash)
		if err != nil {
			return err
		}
		err = newDir.AddChild(metadata, tree)
		return err
	}
}

func Cp(n Node, oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := GetDir(n, oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newDir, err := GetDir(n, newPath)
	if err != nil {
		return err
	}
	for i := 1; i < math.MaxInt64; i++ {
		name := oldName + "(" + strconv.Itoa(i) + ")"
		_, err := GetNode(n, path.Join(newPath, name))
		if err == e.ENoSuchFileOrDir {
			metadata := *oldMetadata
			metadata.Name = name
			return move(newDir, &metadata)
		}
		if err != nil {
			return err
		}
	}
	return e.ErrInvalid
}
