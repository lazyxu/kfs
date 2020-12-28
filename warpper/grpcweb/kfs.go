package grpcweb

import (
	"os"

	"github.com/lazyxu/kfs/kfscore/osmock"

	"github.com/lazyxu/kfs/kfscore/kfscommon"
	"github.com/lazyxu/kfs/kfscore/object"
	"github.com/lazyxu/kfs/kfscore/storage"
)

var obj *object.Obj

func Init(s storage.Storage) {
	obj = object.Init(s)
	_, err := s.GetRef("default")
	if err == nil {
		return
	}
	kfs := osmock.New(&kfscommon.Options{
		UID:       uint32(os.Getuid()),
		GID:       uint32(os.Getgid()),
		DirPerms:  object.S_IFDIR | 0755,
		FilePerms: object.S_IFREG | 0644,
	}, s)
	err = kfs.Storage().UpdateRef("default", "", kfs.Root().Hash())
	if err != nil {
		panic(err)
	}
}
