package core

import (
	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/object"
)

func (kfs *KFS) ReadDir(path string) ([]*object.Metadata, error) {
	n, err := kfs.getNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir.ReadDirAll()
}
