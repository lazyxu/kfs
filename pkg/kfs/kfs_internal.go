package kfs

import (
	"path"
	"strings"

	"github.com/lazyxu/kfs/storage/obj"

	"github.com/lazyxu/kfs/storage/scheduler"

	"github.com/lazyxu/kfs/storage/memory"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/kfs/e"

	"github.com/lazyxu/kfs/kfs/kfscommon"
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
	obj.EmptyDir.Write(kfs.scheduler)
	obj.EmptyFile.Write(kfs.scheduler)
	kfs.root = NewDir(kfs, "")
	kfs.root.Add(obj.NewDirMetadata("demo"), obj.EmptyDir)
	kfs.root.Add(obj.NewFileMetadata("hello"), &obj.File{Reader: strings.NewReader("hello world")})
	kfs.root.Add(obj.NewFileMetadata("index.js"), &obj.File{Reader: strings.NewReader("index")})
	return kfs
}

// GetNode finds the Node by path starting from the root
//
// It is the equivalent of os.Stat - Node contains the os.FileInfo
// interface.
func (kfs *KFS) GetNode(path string) (node Node, err error) {
	defer e.Trace(logrus.Fields{
		"path": path,
	})(func() logrus.Fields {
		return logrus.Fields{
			"err": err,
		}
	})
	return kfs.getNode(path)
}

func (kfs *KFS) getNodeDir(path string) (dir *Dir, err error) {
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

		d, err := obj.ReadDir(kfs.scheduler, dir.Metadata.Hash)
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
					Metadata: *metadata,
				},
				items: make(map[string]Node),
			}
			dir.items[name] = node
		} else {
			node = &File{
				ItemBase: ItemBase{
					kfs:      kfs,
					parent:   dir,
					Metadata: *metadata,
				},
			}
			dir.items[name] = node
		}
	}
	return
}

func (kfs *KFS) getDir(name string) (*Dir, string, error) {
	name = strings.Trim(name, "/")
	parent, leaf := path.Split(name)
	dir, err := kfs.getNodeDir(parent)
	if err != nil {
		return nil, "", err
	}
	return dir, leaf, nil
}
