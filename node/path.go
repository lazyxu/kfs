package node

import (
	"math"
	"path"
	"strconv"
	"strings"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
)

type Mount struct {
	name    string
	root    *Dir
	obj     *object.Obj
	storage storage.Storage
}

func NewMount(root *Dir, name string) *Mount {
	return &Mount{
		name:    name,
		root:    root,
		obj:     root.obj,
		storage: root.storage,
	}
}

func (m *Mount) Commit() string {
	return m.root.Hash
}

func (m *Mount) Obj() *object.Obj {
	return m.obj
}

func (m *Mount) Storage() storage.Storage {
	return m.storage
}

// GetNode finds the Node by path.
func (m *Mount) GetNode(path string) (Node, error) {
	var n Node = m.root
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

		d, err := m.obj.ReadDir(m.storage, dir.Metadata.Hash)
		if err != nil {
			return nil, err
		}
		metadata, err := d.GetNode(name)
		if err != nil {
			return nil, err
		}
		if metadata.IsDir() {
			n = NewDir(m.storage, m.obj, metadata, dir)
			dir.Items[name] = n
		} else {
			n = NewFile(m.storage, m.obj, metadata, dir)
			dir.Items[name] = n
		}
	}
	return n, nil
}

func (m *Mount) GetDir(path string) (*Dir, error) {
	n, err := m.GetNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir, nil
}

func (m *Mount) GetFile(path string) (*File, error) {
	n, err := m.GetNode(path)
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
func (m *Mount) Rename(oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := m.GetDir(oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := m.GetDir(newParent)
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

func (m *Mount) Mv(oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := m.GetDir(oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := m.GetDir(newParent)
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
		newDir, err := m.GetDir(newPath)
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

func (m *Mount) Cp(oldPath, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := m.GetDir(oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newDir, err := m.GetDir(newPath)
	if err != nil {
		return err
	}
	for i := 1; i < math.MaxInt64; i++ {
		name := oldName + "(" + strconv.Itoa(i) + ")"
		_, err := m.GetNode(path.Join(newPath, name))
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
