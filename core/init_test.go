package core

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/lazyxu/kfs/storage/memory"

	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/object"
)

var kfs *KFS

func init() {
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	//storage, _ := fs.New("temp", hashFunc, true, true)
	storage := memory.New(hashFunc)
	serializable := &kfscrypto.GobEncoder{}
	kfs = New(&kfscommon.Options{
		UID:       uint32(os.Getuid()),
		GID:       uint32(os.Getgid()),
		DirPerms:  object.S_IFDIR | 0755,
		FilePerms: object.S_IFREG | 0644,
	}, storage, serializable)
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

func (env *testENV) HasSymlink() bool {
	return false
}

func (env *testENV) HasLink() bool {
	return false
}

var flaky = flag.Bool("flaky", false, "run known-flaky tests too")

func (env *testENV) SkipFlaky(t testing.TB, issue int) {
	t.Helper()
	if !*flaky {
		t.Skipf("skipping known flaky test without the -flaky flag; see golang.org/issue/%d", issue)
	}
}
