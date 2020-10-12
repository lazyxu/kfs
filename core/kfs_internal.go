package core

import (
	"path"
	"strings"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/scheduler"

	"github.com/lazyxu/kfs/storage/memory"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/core/kfscommon"
)

type KFS struct {
	root      *Dir
	scheduler *scheduler.Scheduler
	Opt       *kfscommon.Options
}

func New(opt *kfscommon.Options) *KFS {
	kfs := &KFS{
		Opt:       opt,
		scheduler: scheduler.New(memory.New()),
	}
	object.EmptyDir.Write(kfs.scheduler)
	object.EmptyFile.Write(kfs.scheduler)
	kfs.root = NewDir(kfs, "", object.DefaultDirMode)
	kfs.root.add(object.NewDirMetadata("demo", object.DefaultDirMode), object.EmptyDir)
	kfs.root.add(object.NewFileMetadata("hello"), &object.Blob{Reader: strings.NewReader("hello world")})
	kfs.root.add(object.NewFileMetadata("index.js"), &object.Blob{Reader: strings.NewReader("index")})
	return kfs
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
			return nil, e.ErrNotExist
		}
		node, ok = dir.items[name]
		if ok {
			continue
		}

		d, err := object.ReadDir(kfs.scheduler, dir.Metadata.Hash)
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
	name = strings.Trim(name, "/")
	parent, leaf := path.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return nil, "", err
	}
	return dir, leaf, nil
}
