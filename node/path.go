package node

import (
	"strings"

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

//func Copy(n Node, srcList []string, dst string) error {
//	obj := n.Obj()
//	storage := n.Storage()
//	if len(srcList) == 0 {
//		return nil
//	}
//	if len(srcList) == 1 {
//		src := srcList[0]
//		srcNode, err := GetNode(n, src)
//		if err != nil {
//			return err
//		}
//		dstNode, err := GetNode(n, dst)
//		if err != nil {
//			return err
//		}
//		if srcNode == dstNode {
//			return nil
//		}
//		if dir, ok := dstNode.(*Dir); ok {
//			metadata, err := dir.GetChild(srcNode.Name())
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return n, nil
//}
