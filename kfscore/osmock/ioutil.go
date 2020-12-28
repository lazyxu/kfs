package osmock

import (
	"github.com/lazyxu/kfs/kfscore/e"
	"github.com/lazyxu/kfs/kfscore/node"
	"github.com/lazyxu/kfs/kfscore/object"
)

func (kfs *KFS) ReadDir(path string) ([]*object.Metadata, error) {
	n, err := kfs.GetNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*node.Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir.ReadDirAll()
}
