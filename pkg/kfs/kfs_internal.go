package kfs

import (
	"path"
	"strings"

	"github.com/lazyxu/kfs/storage/scheduler"

	"github.com/lazyxu/kfs/storage/memory"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/kfs/e"
	"github.com/lazyxu/kfs/node"

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
	kfs.root = NewDir(kfs, "")
	kfs.root.Add("demo", NewDir(kfs, "demo"))
	hello, _ := NewFile(kfs, "hello")
	hello.SetContent([]byte("hello world"), 0)
	kfs.root.Add("hello", hello)
	index, _ := NewFile(kfs, "index.js")
	index.SetContent([]byte("index"), 0)
	kfs.root.Add("index.js", index)
	return kfs
}

// GetNode finds the Node by path starting from the root
//
// It is the equivalent of os.Stat - Node contains the os.FileInfo
// interface.
func (kfs *KFS) GetNode(path string) (node node.Node, err error) {
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

func (kfs *KFS) getNodeFile(path string) (file *File, err error) {
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

func (kfs *KFS) getNode(path string) (node node.Node, err error) {
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
		item, err := dir.Stat(name)
		if err != nil {
			return nil, err
		}
		node, err = item.Node()
		if err != nil {
			return nil, err
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
