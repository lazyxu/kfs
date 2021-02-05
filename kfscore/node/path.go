package node

import (
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/lazyxu/kfs/kfscore/object"

	"github.com/lazyxu/kfs/kfscore/e"
)

type Mount struct {
	name    string
	head    string
	root    *Dir
	obj     *object.Obj
	storage storage.Storage
}

func NewMount(name string, s storage.Storage) (*Mount, error) {
	obj := object.Init(s)
	head, err := s.GetRef(name)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	root := NewDir(s, obj,
		obj.NewDirMetadata(name, object.DefaultDirMode).Builder().Hash(head).Build(), nil)
	return &Mount{
		name:    name,
		head:    head,
		root:    root,
		obj:     obj,
		storage: s,
	}, nil
}

func (m *Mount) Commit() error {
	err := m.storage.UpdateRef(m.name, m.head, m.root.Hash())
	if err != nil {
		return err
	}
	m.head = m.root.Hash()
	return nil
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

		d, err := m.obj.ReadTree(dir.metadata.Hash())
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
		return move(newDir, oldMetadata, newName)
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
		return move(newDir, oldMetadata, newName)
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
	if newName == "" {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		newDir, err := m.GetDir(newPath)
		if err != nil {
			return err
		}
		return move(newDir, oldMetadata, oldName)
	}
	newMetadata, err := newDir.GetChild(newName)
	if err == e.ErrNotExist {
		err := oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		return move(newDir, oldMetadata, newName)
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
		return move(newDir, oldMetadata, newName)
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
		return move(newDir, oldMetadata, oldName)
	}
	return nil
}

func Mv(mSrc *Mount, oldPath string, mDst *Mount, newPath string) error {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := mSrc.GetDir(oldParent)
	if err != nil {
		return err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return err
	}
	newParent, newName := path.Split(newPath)
	newDir, err := mDst.GetDir(newParent)
	if err != nil {
		return err
	}
	if newName == "" {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		newDir, err := mDst.GetDir(newPath)
		if err != nil {
			return err
		}
		return move(newDir, oldMetadata, oldName)
	}
	newMetadata, err := newDir.GetChild(newName)
	if err == e.ErrNotExist {
		err := oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		return move(newDir, oldMetadata, newName)
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
		return move(newDir, oldMetadata, newName)
	}
	if newMetadata.IsDir() {
		err = oldDir.RemoveChild(oldName, true)
		if err != nil {
			return err
		}
		newDir, err := mDst.GetDir(newPath)
		if err != nil {
			return err
		}
		return move(newDir, oldMetadata, oldName)
	}
	return nil
}

func move(newDir *Dir, oldMetadata *object.Metadata, name string) error {
	metadata := oldMetadata.Builder().Name(name).Build()
	if metadata.IsFile() {
		err := newDir.AddChild(metadata)
		return err
	} else {
		err := newDir.AddChild(metadata)
		return err
	}
}

func Cp(mSrc *Mount, oldPath string, mDst *Mount, newPath string) (string, error) {
	oldParent, oldName := path.Split(oldPath)
	oldDir, err := mSrc.GetDir(oldParent)
	if err != nil {
		return "", err
	}
	oldMetadata, err := oldDir.GetChild(oldName)
	if err != nil {
		return "", err
	}
	newDir, err := mDst.GetDir(newPath)
	if err != nil {
		return "", err
	}
	return mDst.tryName(newPath, oldName, func(name string) error {
		return move(newDir, oldMetadata, name)
	})
}

func (m *Mount) tryName(p string, baseName string, fn func(name string) error) (string, error) {
	name := baseName
	_, err := m.GetNode(path.Join(p, name))
	if err == e.ENoSuchFileOrDir {
		return name, fn(name)
	}
	if err != nil {
		return "", err
	}
	for i := 1; i < math.MaxInt64; i++ {
		name = baseName + "(" + strconv.Itoa(i) + ")"
		_, err := m.GetNode(path.Join(p, name))
		if err == e.ENoSuchFileOrDir {
			return name, fn(name)
		}
		if err != nil {
			return "", err
		}
	}
	return "", e.ErrInvalid
}

func (m *Mount) NewFile(p string) (string, error) {
	dir, err := m.GetDir(p)
	if err != nil {
		return "", err
	}
	return m.tryName(p, "未命名文件", func(name string) error {
		return dir.AddChild(m.Obj().NewFileMetadata(name, object.DefaultFileMode))
	})
}

func (m *Mount) NewDir(p string) (string, error) {
	dir, err := m.GetDir(p)
	if err != nil {
		return "", err
	}
	return m.tryName(p, "未命名文件夹", func(name string) error {
		return dir.AddChild(m.Obj().NewDirMetadata(name, object.DefaultDirMode))
	})
}
