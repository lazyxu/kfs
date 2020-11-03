package core

import (
	"testing"

	"github.com/lazyxu/kfs/object"
)

func init() {
	kfs.Mkdir("/etc", object.DefaultDirMode)
	kfs.Create("/etc/group")
	kfs.Create("/etc/hosts")
	kfs.Create("/etc/passwd")
	kfs.Mkdir("/tmp", object.DefaultDirMode)
}

var testenv testENV

type testENV struct {
}

func (env *testENV) MustHaveSymlink(t testing.TB) {
	t.Skipf("skipping test: cannot make symlinks")
}
