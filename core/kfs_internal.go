package core

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/core/kfscommon"
)

type KFS struct {
	baseObject *object.BaseObject
	root       *Dir
	storage    storage.Storage
	Opt        *kfscommon.Options
	pwd        string
}

var defaultBinaries = []string{
	"cat",
	"chmod",
	"cp",
	"date",
	"dd",
	"df",
	"hostname",
	"link",
	"ln",
	"ls",
	"mkdir",
	"mv",
	"ps",
	"pwd",
	"rm",
	"rmdir",
	"sync",
	"unlink",
}

const DevNull = "/dev/null"

func New(opt *kfscommon.Options, s storage.Storage,
	hashFunc func() kfscrypto.Hash, serializable kfscrypto.Serializable) *KFS {
	baseObject := object.Init(hashFunc, serializable)
	kfs := &KFS{
		Opt:        opt,
		storage:    s,
		pwd:        "/tmp",
		baseObject: baseObject,
	}
	baseObject.EmptyDir.Write(kfs.storage)
	baseObject.EmptyFile.Write(kfs.storage)
	kfs.root = NewDir(kfs, "", object.DefaultDirMode)
	err := kfs.root.add(baseObject.NewDirMetadata("demo", object.DefaultDirMode), baseObject.EmptyDir)
	if err != nil {
		panic(err)
	}
	err = kfs.root.add(baseObject.NewFileMetadata("hello"), &object.Blob{Reader: strings.NewReader("hello world")})
	if err != nil {
		panic(err)
	}
	err = kfs.root.add(baseObject.NewFileMetadata("index.js"), &object.Blob{Reader: strings.NewReader("index")})
	if err != nil {
		panic(err)
	}
	err = kfs.MkdirAll("/home/test", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	err = kfs.Mkdir("/bin", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	for _, b := range defaultBinaries {
		_, err = kfs.Create(path.Join("/bin", b))
		if err != nil {
			panic(err)
		}
	}
	err = kfs.Mkdir("/tmp", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	return kfs
}

func (kfs *KFS) Getwd() (dir string, err error) {
	return kfs.pwd, nil
}

func (kfs *KFS) GetDir(path string) (dir *Dir, err error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir, nil
}

func (kfs *KFS) GetFile(path string) (*File, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return nil, err
	}
	file, ok := n.(*File)
	if !ok {
		return nil, e.ENotFile
	}
	return file, nil
}

// getNode finds the Node by path starting from the root
func (kfs *KFS) getNode(path string) (node Node, err error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(kfs.pwd, path)
	}
	path = strings.Trim(path, "/")
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
			return nil, e.ENotDir
		}
		node, ok = dir.items[name]
		if ok {
			continue
		}

		d, err := kfs.baseObject.ReadDir(kfs.storage, dir.Metadata.Hash)
		if err != nil {
			return nil, err
		}
		metadata, err := d.GetNode(name)
		if err != nil {
			return nil, err
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
	return
}

func (kfs *KFS) getDirAndLeaf(name string) (*Dir, string, error) {
	parent, leaf := filepath.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return nil, "", err
	}
	return dir, leaf, nil
}
