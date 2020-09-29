package kfs

import (
	"fmt"
	"path"
	"strings"

	"github.com/lazyxu/kfs/storage/scheduler"

	"github.com/lazyxu/kfs/storage/memory"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/kfs/e"

	"github.com/lazyxu/kfs/kfs/kfscommon"
)

type KFS struct {
	itemRoot  *ItemDir
	scheduler *scheduler.Scheduler
	Opt       *kfscommon.Options
}

func New(opt *kfscommon.Options) *KFS {
	kfs := &KFS{
		Opt:       opt,
		scheduler: scheduler.New(memory.New()),
	}
	kfs.itemRoot = newDir(kfs, "")
	kfs.itemRoot.Add(newDir(kfs, "demo"))
	hello := newFile(kfs, "hello")
	fmt.Println("Hello done")
	hello.SetContent([]byte("hello world"), 0)
	fmt.Println("Hello done")
	kfs.itemRoot.Add(hello)
	index := newFile(kfs, "index.js")
	index.SetContent([]byte("index"), 0)
	kfs.itemRoot.Add(index)
	return kfs
}

// GetNode finds the Node by path starting from the root
//
// It is the equivalent of os.Stat - Node contains the os.FileInfo
// interface.
func (kfs *KFS) GetNode(path string) (node Item, err error) {
	defer e.Trace(logrus.Fields{
		"path": path,
	})(func() logrus.Fields {
		return logrus.Fields{
			"err": err,
		}
	})
	return kfs.getNode(path)
}

func (kfs *KFS) getNodeDir(path string) (dir *ItemDir, err error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*ItemDir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir, nil
}

func (kfs *KFS) GetFile(path string) (file *ItemFile, err error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return nil, err
	}
	file, ok := n.(*ItemFile)
	if !ok {
		return nil, e.ENotFile
	}
	return file, nil
}

func (kfs *KFS) getNode(path string) (node Item, err error) {
	path = strings.Trim(path, "/")
	node = kfs.itemRoot
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
		dir, ok := node.(*ItemDir)
		if !ok {
			// We need to look in a directory, but found a file
			return nil, e.ErrNotExist
		}
		node, err = dir.GetNode(name)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (kfs *KFS) getDir(name string) (*ItemDir, string, error) {
	name = strings.Trim(name, "/")
	parent, leaf := path.Split(name)
	dir, err := kfs.getNodeDir(parent)
	if err != nil {
		return nil, "", err
	}
	return dir, leaf, nil
}
