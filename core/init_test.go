package core

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/lazyxu/kfs/object"
)

func init() {
	kfs.Mkdir("/etc", object.DefaultDirMode)
	group, _ := kfs.Create("/etc/group")
	group.WriteAt([]byte(strings.Repeat("x", 1000)), 0)
	group.Close()
	kfs.Create("/etc/hosts")
	kfs.Create("/etc/passwd")
	kfs.Mkdir("/tmp", object.DefaultDirMode)
	err := filepath.Walk(path.Join(runtime.GOROOT(), "src/os"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		err = kfs.MkdirAll(path, kfs.Opt.DirPerms)
		if err != nil {
			return err
		}
		dst, err := kfs.Create(path)
		if err != nil {
			return err
		}
		defer dst.Close()
		srcBuf := bufio.NewReader(src)
		dstBuf := bufio.NewWriter(dst)
		srcBuf.WriteTo(dstBuf)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

var testenv testENV

type testENV struct {
}

func (env *testENV) MustHaveSymlink(t testing.TB) {
	t.Skipf("skipping test: cannot make symlinks")
}

func (env *testENV) MustHaveLink(t testing.TB) {
	t.Skipf("skipping test: cannot make hard links")
}

func (env *testENV) MustHaveExec(t testing.TB) {
	t.Skipf("skipping test: cannot exec")
}
