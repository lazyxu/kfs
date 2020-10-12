package core

import (
	"os"
	"testing"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/object"
)

var kfs = New(&kfscommon.Options{
	UID:       uint32(os.Getuid()),
	GID:       uint32(os.Getgid()),
	DirPerms:  object.S_IFDIR | 0755,
	FilePerms: object.S_IFREG | 0644,
})

func init() {
	etc := NewDir(kfs, "etc", object.DefaultDirMode)
	etc.add(object.NewFileMetadata("group"), object.EmptyFile)
	etc.add(object.NewFileMetadata("hosts"), object.EmptyFile)
	etc.add(object.NewFileMetadata("passwd"), object.EmptyFile)
	tree, _ := etc.load()
	kfs.root.add(etc.Metadata, tree)
	kfs.root.add(object.NewDirMetadata("tmp", object.DefaultDirMode), object.EmptyDir)
}

var testenv testENV

type testENV struct {
}

func (env *testENV) MustHaveSymlink(t testing.TB) {
	t.Skipf("skipping test: cannot make symlinks")
}

func Open(name string) (*File, error) {
	return kfs.Open(name)
}

func Create(name string) (*File, error) {
	return kfs.Create(name)
}

func Stat(name string) (os.FileInfo, error) {
	return kfs.Stat(name)
}

// Symlink creates newname as a symbolic link to oldname.
func Symlink(oldname, newname string) error {
	return e.ENotImpl
}
