// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	osexec "os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"
)

var dot = []string{
	"dir.go",
	"file.go",
	"node.go",
	"os_test.go",
	"mock_test.go",
	"init_test.go",
	"stat.go",
	"handle.go",
	"handle_dir.go",
	"handle_write.go",
	"handle_read.go",
}

type sysDir struct {
	name  string
	files []string
}

var sysdir = func() *sysDir {
	/*
		align with os/os_test.go




















































	*/
	return &sysDir{
		"/etc",
		[]string{
			"group",
			"hosts",
			"passwd",
		},
	}
}()

func size(name string, t *testing.T) int64 {
	file, err := Open(name)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer file.Close()
	var buf [100]byte
	len := 0
	for {
		n, e := file.Read(buf[0:])
		len += n
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal("read failed:", e)
		}
	}
	return int64(len)
}

func equal(name1, name2 string) bool {
	/*
		align with os/os_test.go



	*/
	return name1 == name2
}

// localTmp returns a local temporary directory not on NFS.
func localTmp() string {
	/*
		align with os/os_test.go






	*/
	return "/tmp"
}

func newFile(testName string, t *testing.T) (fh Handle) {
	fh, err := Create(path.Join(localTmp(), "_Go_"+testName))
	if err != nil {
		t.Fatalf("TempFile %s: %s", testName, err)
	}
	return
}

func newDir(testName string, t *testing.T) (name string) {
	name, err := TempDir(localTmp(), "_Go_"+testName)
	if err != nil {
		t.Fatalf("TempDir %s: %s", testName, err)
	}
	return
}

var sfdir = sysdir.name
var sfname = sysdir.files[0]

func TestStat(t *testing.T) {
	path := sfdir + "/" + sfname
	dir, err := Stat(path)
	if err != nil {
		t.Fatal("stat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

func TestStatError(t *testing.T) {
	defer chtmpdir(t)()

	path := "no-such-file"

	fi, err := Stat(path)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if fi != nil {
		t.Errorf("got %v, want nil", fi)
	}
	if perr, ok := err.(*PathError); !ok {
		t.Errorf("got %T, want %T", err, perr)
	}

	testenv.MustHaveSymlink(t)

	link := "symlink"
	err = Symlink(path, link)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = Stat(link)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if fi != nil {
		t.Errorf("got %v, want nil", fi)
	}
	if perr, ok := err.(*PathError); !ok {
		t.Errorf("got %T, want %T", err, perr)
	}
}

func TestFstat(t *testing.T) {
	path := sfdir + "/" + sfname
	file, err1 := Open(path)
	if err1 != nil {
		t.Fatal("open failed:", err1)
	}
	defer file.Close()
	dir, err2 := file.Stat()
	if err2 != nil {
		t.Fatal("fstat failed:", err2)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

func TestLstat(t *testing.T) {
	path := sfdir + "/" + sfname
	dir, err := Lstat(path)
	if err != nil {
		t.Fatal("lstat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

// Read with length 0 should not return EOF.
func TestRead0(t *testing.T) {
	path := sfdir + "/" + sfname
	f, err := Open(path)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer f.Close()

	b := make([]byte, 0)
	n, err := f.Read(b)
	if n != 0 || err != nil {
		t.Errorf("Read(0) = %d, %v, want 0, nil", n, err)
	}
	b = make([]byte, 100)
	n, err = f.Read(b)
	if n <= 0 || err != nil {
		t.Errorf("Read(100) = %d, %v, want >0, nil", n, err)
	}
}

// Reading a closed file should return ErrClosed error
func TestReadClosed(t *testing.T) {
	path := sfdir + "/" + sfname
	file, err := Open(path)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	file.Close() // close immediately

	b := make([]byte, 100)
	_, err = file.Read(b)

	e, ok := err.(*PathError)
	if !ok {
		t.Fatalf("Read: %T(%v), want PathError", e, e)
	}

	if e.Err != ErrClosed {
		t.Errorf("Read: %v, want PathError(ErrClosed)", e)
	}
}

func testReaddirnames(dir string, contents []string, t *testing.T) {
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	s, err2 := file.Readdirnames(-1)
	if err2 != nil {
		t.Fatalf("readdirnames %q failed: %v", dir, err2)
	}
	for _, m := range contents {
		found := false
		for _, n := range s {
			if n == "." || n == ".." {
				t.Errorf("got %s in directory", n)
			}
			if equal(m, n) {
				if found {
					t.Error("present twice:", m)
				}
				found = true
			}
		}
		if !found {
			t.Error("could not find", m)
		}
	}
}

func testReaddir(dir string, contents []string, t *testing.T) {
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	s, err2 := file.Readdir(-1)
	if err2 != nil {
		t.Fatalf("readdir %q failed: %v", dir, err2)
	}
	for _, m := range contents {
		found := false
		for _, n := range s {
			if equal(m, n.Name()) {
				if found {
					t.Error("present twice:", m)
				}
				found = true
			}
		}
		if !found {
			t.Error("could not find", m)
		}
	}
}

func TestReaddirnames(t *testing.T) {
	//testReaddirnames(".", dot, t)
	testReaddirnames(sysdir.name, sysdir.files, t)
}

func TestReaddir(t *testing.T) {
	//testReaddir(".", dot, t)
	testReaddir(sysdir.name, sysdir.files, t)
}

func benchmarkReaddirname(path string, b *testing.B) {
	var nentries int
	for i := 0; i < b.N; i++ {
		f, err := Open(path)
		if err != nil {
			b.Fatalf("open %q failed: %v", path, err)
		}
		ns, err := f.Readdirnames(-1)
		f.Close()
		if err != nil {
			b.Fatalf("readdirnames %q failed: %v", path, err)
		}
		nentries = len(ns)
	}
	b.Logf("benchmarkReaddirname %q: %d entries", path, nentries)
}

func benchmarkReaddir(path string, b *testing.B) {
	var nentries int
	for i := 0; i < b.N; i++ {
		f, err := Open(path)
		if err != nil {
			b.Fatalf("open %q failed: %v", path, err)
		}
		fs, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			b.Fatalf("readdir %q failed: %v", path, err)
		}
		nentries = len(fs)
	}
	b.Logf("benchmarkReaddir %q: %d entries", path, nentries)
}

func BenchmarkReaddirname(b *testing.B) {
	benchmarkReaddirname("/etc", b)
}

func BenchmarkReaddir(b *testing.B) {
	benchmarkReaddir("[root]", b)
}

func benchmarkStat(b *testing.B, path string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Stat(path)
		if err != nil {
			b.Fatalf("Stat(%q) failed: %v", path, err)
		}
	}
}

func benchmarkLstat(b *testing.B, path string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Lstat(path)
		if err != nil {
			b.Fatalf("Lstat(%q) failed: %v", path, err)
		}
	}
}

func BenchmarkStatDot(b *testing.B) {
	benchmarkStat(b, ".")
}

func BenchmarkStatFile(b *testing.B) {
	benchmarkStat(b, filepath.Join(runtime.GOROOT(), "src/os/os_test.go"))
}

func BenchmarkStatDir(b *testing.B) {
	benchmarkStat(b, filepath.Join(runtime.GOROOT(), "src/os"))
}

func BenchmarkLstatDot(b *testing.B) {
	benchmarkLstat(b, ".")
}

func BenchmarkLstatFile(b *testing.B) {
	benchmarkLstat(b, filepath.Join(runtime.GOROOT(), "src/os/os_test.go"))
}

func BenchmarkLstatDir(b *testing.B) {
	benchmarkLstat(b, filepath.Join(runtime.GOROOT(), "src/os"))
}

// Read the directory one entry at a time.
func smallReaddirnames(file Handle, length int, t *testing.T) []string {
	names := make([]string, length)
	count := 0
	for {
		d, err := file.Readdirnames(1)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("readdirnames %q failed: %v", file.Name(), err)
		}
		if len(d) == 0 {
			t.Fatalf("readdirnames %q returned empty slice and no error", file.Name())
		}
		names[count] = d[0]
		count++
	}
	return names[0:count]
}

// Check that reading a directory one entry at a time gives the same result
// as reading it all at once.
func TestReaddirnamesOneAtATime(t *testing.T) {
	// big directory that doesn't change often.
	dir := "/bin"
	/*
		align with os/os_test.go














	*/
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	all, err1 := file.Readdirnames(-1)
	if err1 != nil {
		t.Fatalf("readdirnames %q failed: %v", dir, err1)
	}
	file1, err2 := Open(dir)
	if err2 != nil {
		t.Fatalf("open %q failed: %v", dir, err2)
	}
	defer file1.Close()
	small := smallReaddirnames(file1, len(all)+100, t) // +100 in case we screw up
	if len(small) < len(all) {
		t.Fatalf("len(small) is %d, less than %d", len(small), len(all))
	}
	for i, n := range all {
		if small[i] != n {
			t.Errorf("small read %q mismatch: %v", small[i], n)
		}
	}
}

func TestReaddirNValues(t *testing.T) {
	if testing.Short() {
		t.Skip("test.short; skipping")
	}
	dir, err := TempDir("", "")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer RemoveAll(dir)
	for i := 1; i <= 105; i++ {
		f, err := Create(filepath.Join(dir, fmt.Sprintf("%d", i)))
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		f.Write([]byte(strings.Repeat("X", i)))
		f.Close()
	}

	var d Handle
	openDir := func() {
		var err error
		d, err = Open(dir)
		if err != nil {
			t.Fatalf("Open directory: %v", err)
		}
	}

	readDirExpect := func(n, want int, wantErr error) {
		fi, err := d.Readdir(n)
		if err != wantErr {
			t.Fatalf("Readdir of %d got error %v, want %v", n, err, wantErr)
		}
		if g, e := len(fi), want; g != e {
			t.Errorf("Readdir of %d got %d files, want %d", n, g, e)
		}
	}

	readDirNamesExpect := func(n, want int, wantErr error) {
		fi, err := d.Readdirnames(n)
		if err != wantErr {
			t.Fatalf("Readdirnames of %d got error %v, want %v", n, err, wantErr)
		}
		if g, e := len(fi), want; g != e {
			t.Errorf("Readdirnames of %d got %d files, want %d", n, g, e)
		}
	}

	for _, fn := range []func(int, int, error){readDirExpect, readDirNamesExpect} {
		// Test the slurp case
		openDir()
		fn(0, 105, nil)
		fn(0, 0, nil)
		d.Close()

		// Slurp with -1 instead
		openDir()
		fn(-1, 105, nil)
		fn(-2, 0, nil)
		fn(0, 0, nil)
		d.Close()

		// Test the bounded case
		openDir()
		fn(1, 1, nil)
		fn(2, 2, nil)
		fn(105, 102, nil) // and tests buffer >100 case
		fn(3, 0, io.EOF)
		d.Close()
	}
}

func touch(t *testing.T, name string) {
	f, err := Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestReaddirStatFailures(t *testing.T) {
	// KFS already do this correctly,
	// but are structured with different syscalls such
	// that they don't use Lstat, so the hook below for
	// testing it wouldn't work.
	/*
		align with os/os_test.go
	*/
	t.Skipf("skipping test")
	dir, err := TempDir("", "")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer RemoveAll(dir)
	touch(t, filepath.Join(dir, "good1"))
	touch(t, filepath.Join(dir, "x")) // will disappear or have an error
	touch(t, filepath.Join(dir, "good2"))
	defer func() {
		*LstatP = Lstat
	}()
	var xerr error // error to return for x
	*LstatP = func(path string) (FileInfo, error) {
		if xerr != nil && strings.HasSuffix(path, "x") {
			return nil, xerr
		}
		return Lstat(path)
	}
	readDir := func() ([]FileInfo, error) {
		d, err := Open(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer d.Close()
		return d.Readdir(-1)
	}
	mustReadDir := func(testName string) []os.FileInfo {
		fis, err := readDir()
		if err != nil {
			t.Fatalf("%s: Readdir: %v", testName, err)
		}
		return fis
	}
	names := func(fis []os.FileInfo) []string {
		s := make([]string, len(fis))
		for i, fi := range fis {
			s[i] = fi.Name()
		}
		sort.Strings(s)
		return s
	}

	if got, want := names(mustReadDir("initial readdir")),
		[]string{"good1", "good2", "x"}; !reflect.DeepEqual(got, want) {
		t.Errorf("initial readdir got %q; want %q", got, want)
	}

	xerr = ErrNotExist
	if got, want := names(mustReadDir("with x disappearing")),
		[]string{"good1", "good2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("with x disappearing, got %q; want %q", got, want)
	}

	xerr = errors.New("some real error")
	if _, err := readDir(); err != xerr {
		t.Errorf("with a non-ErrNotExist error, got error %v; want %v", err, xerr)
	}
}

// Readdir on a regular file should fail.
func TestReaddirOfFile(t *testing.T) {
	f, err := TempFile("", "_Go_ReaddirOfFile")
	if err != nil {
		t.Fatal(err)
	}
	defer Remove(f.Name())
	f.Write([]byte("foo"))
	f.Close()
	reg, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer reg.Close()

	names, err := reg.Readdirnames(-1)
	if err == nil {
		t.Error("Readdirnames succeeded; want non-nil error")
	}
	if len(names) > 0 {
		t.Errorf("unexpected dir names in regular file: %q", names)
	}
}

func TestHardLink(t *testing.T) {
	testenv.MustHaveLink(t)

	defer chtmpdir(t)()
	from, to := "hardlinktestfrom", "hardlinktestto"
	file, err := Create(to)
	if err != nil {
		t.Fatalf("open %q failed: %v", to, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("close %q failed: %v", to, err)
	}
	err = Link(to, from)
	if err != nil {
		t.Fatalf("link %q, %q failed: %v", to, from, err)
	}

	none := "hardlinktestnone"
	err = Link(none, none)
	// Check the returned error is well-formed.
	if lerr, ok := err.(*LinkError); !ok || lerr.Error() == "" {
		t.Errorf("link %q, %q failed to return a valid error", none, none)
	}

	tostat, err := Stat(to)
	if err != nil {
		t.Fatalf("stat %q failed: %v", to, err)
	}
	fromstat, err := Stat(from)
	if err != nil {
		t.Fatalf("stat %q failed: %v", from, err)
	}
	if !SameFile(tostat, fromstat) {
		t.Errorf("link %q, %q did not create hard link", to, from)
	}
	// We should not be able to perform the same Link() a second time
	err = Link(to, from)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "link" {
			t.Errorf("Link(%q, %q) err.Op = %q; want %q", to, from, err.Op, "link")
		}
		if err.Old != to {
			t.Errorf("Link(%q, %q) err.Old = %q; want %q", to, from, err.Old, to)
		}
		if err.New != from {
			t.Errorf("Link(%q, %q) err.New = %q; want %q", to, from, err.New, from)
		}
		if !IsExist(err.Err) {
			t.Errorf("Link(%q, %q) err.Err = %q; want %q", to, from, err.Err, "file exists error")
		}
	case nil:
		t.Errorf("link %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("link %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

// chtmpdir changes the working directory to a new temporary directory and
// provides a cleanup function.
func chtmpdir(t *testing.T) func() {
	oldwd, err := Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := TempDir("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		RemoveAll(d)
	}
}

func TestSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)

	defer chtmpdir(t)()
	from, to := "symlinktestfrom", "symlinktestto"
	file, err := Create(to)
	if err != nil {
		t.Fatalf("Create(%q) failed: %v", to, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("Close(%q) failed: %v", to, err)
	}
	err = Symlink(to, from)
	if err != nil {
		t.Fatalf("Symlink(%q, %q) failed: %v", to, from, err)
	}
	tostat, err := Lstat(to)
	if err != nil {
		t.Fatalf("Lstat(%q) failed: %v", to, err)
	}
	if tostat.Mode()&ModeSymlink != 0 {
		t.Fatalf("Lstat(%q).Mode()&ModeSymlink = %v, want 0", to, tostat.Mode()&ModeSymlink)
	}
	fromstat, err := Stat(from)
	if err != nil {
		t.Fatalf("Stat(%q) failed: %v", from, err)
	}
	if !SameFile(tostat, fromstat) {
		t.Errorf("Symlink(%q, %q) did not create symlink", to, from)
	}
	fromstat, err = Lstat(from)
	if err != nil {
		t.Fatalf("Lstat(%q) failed: %v", from, err)
	}
	if fromstat.Mode()&ModeSymlink == 0 {
		t.Fatalf("Lstat(%q).Mode()&ModeSymlink = 0, want %v", from, ModeSymlink)
	}
	fromstat, err = Stat(from)
	if err != nil {
		t.Fatalf("Stat(%q) failed: %v", from, err)
	}
	if fromstat.Name() != from {
		t.Errorf("Stat(%q).Name() = %q, want %q", from, fromstat.Name(), from)
	}
	if fromstat.Mode()&ModeSymlink != 0 {
		t.Fatalf("Stat(%q).Mode()&ModeSymlink = %v, want 0", from, fromstat.Mode()&ModeSymlink)
	}
	s, err := Readlink(from)
	if err != nil {
		t.Fatalf("Readlink(%q) failed: %v", from, err)
	}
	if s != to {
		t.Fatalf("Readlink(%q) = %q, want %q", from, s, to)
	}
	file, err = Open(from)
	if err != nil {
		t.Fatalf("Open(%q) failed: %v", from, err)
	}
	file.Close()
}

func TestLongSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)

	defer chtmpdir(t)()
	s := "0123456789abcdef"
	// Long, but not too long: a common limit is 255.
	s = s + s + s + s + s + s + s + s + s + s + s + s + s + s + s
	from := "longsymlinktestfrom"
	err := Symlink(s, from)
	if err != nil {
		t.Fatalf("symlink %q, %q failed: %v", s, from, err)
	}
	r, err := Readlink(from)
	if err != nil {
		t.Fatalf("readlink %q failed: %v", from, err)
	}
	if r != s {
		t.Fatalf("after symlink %q != %q", r, s)
	}
}

func TestRename(t *testing.T) {
	defer chtmpdir(t)()
	from, to := "renamefrom", "renameto"

	file, err := Create(from)
	if err != nil {
		t.Fatalf("open %q failed: %v", from, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("close %q failed: %v", from, err)
	}
	err = Rename(from, to)
	if err != nil {
		t.Fatalf("rename %q, %q failed: %v", to, from, err)
	}
	_, err = Stat(to)
	if err != nil {
		t.Errorf("stat %q failed: %v", to, err)
	}
}

func TestRenameOverwriteDest(t *testing.T) {
	defer chtmpdir(t)()
	from, to := "renamefrom", "renameto"

	toData := []byte("to")
	fromData := []byte("from")

	err := WriteFile(to, toData, 0777)
	if err != nil {
		t.Fatalf("write file %q failed: %v", to, err)
	}

	err = WriteFile(from, fromData, 0777)
	if err != nil {
		t.Fatalf("write file %q failed: %v", from, err)
	}
	err = Rename(from, to)
	if err != nil {
		t.Fatalf("rename %q, %q failed: %v", to, from, err)
	}

	_, err = Stat(from)
	if err == nil {
		t.Errorf("from file %q still exists", from)
	}
	if err != nil && !IsNotExist(err) {
		t.Fatalf("stat from: %v", err)
	}
	toFi, err := Stat(to)
	if err != nil {
		t.Fatalf("stat %q failed: %v", to, err)
	}
	if toFi.Size() != int64(len(fromData)) {
		t.Errorf(`"to" size = %d; want %d (old "from" size)`, toFi.Size(), len(fromData))
	}
}

func TestRenameFailed(t *testing.T) {
	defer chtmpdir(t)()
	from, to := "renamefrom", "renameto"

	err := Rename(from, to)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "rename" {
			t.Errorf("rename %q, %q: err.Op: want %q, got %q", from, to, "rename", err.Op)
		}
		if err.Old != from {
			t.Errorf("rename %q, %q: err.Old: want %q, got %q", from, to, from, err.Old)
		}
		if err.New != to {
			t.Errorf("rename %q, %q: err.New: want %q, got %q", from, to, to, err.New)
		}
	case nil:
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("rename %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

func TestRenameNotExisting(t *testing.T) {
	defer chtmpdir(t)()
	from, to := "doesnt-exist", "dest"

	Mkdir(to, 0777)

	if err := Rename(from, to); !IsNotExist(err) {
		t.Errorf("Rename(%q, %q) = %v; want an IsNotExist error", from, to, err)
	}
}

func TestRenameToDirFailed(t *testing.T) {
	defer chtmpdir(t)()
	from, to := "renamefrom", "renameto"

	Mkdir(from, 0777)
	Mkdir(to, 0777)

	err := Rename(from, to)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "rename" {
			t.Errorf("rename %q, %q: err.Op: want %q, got %q", from, to, "rename", err.Op)
		}
		if err.Old != from {
			t.Errorf("rename %q, %q: err.Old: want %q, got %q", from, to, from, err.Old)
		}
		if err.New != to {
			t.Errorf("rename %q, %q: err.New: want %q, got %q", from, to, to, err.New)
		}
	case nil:
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("rename %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

func TestRenameCaseDifference(pt *testing.T) {
	from, to := "renameFROM", "RENAMEfrom"
	tests := []struct {
		name   string
		create func() error
	}{
		{"dir", func() error {
			return Mkdir(from, 0777)
		}},
		{"file", func() error {
			fd, err := Create(from)
			if err != nil {
				return err
			}
			return fd.Close()
		}},
	}

	for _, test := range tests {
		pt.Run(test.name, func(t *testing.T) {
			defer chtmpdir(t)()

			if err := test.create(); err != nil {
				t.Fatalf("failed to create test file: %s", err)
			}

			if _, err := Stat(to); err != nil {
				// Sanity check that the underlying filesystem is not case sensitive.
				if IsNotExist(err) {
					t.Skipf("case sensitive filesystem")
				}
				t.Fatalf("stat %q, got: %q", to, err)
			}

			if err := Rename(from, to); err != nil {
				t.Fatalf("unexpected error when renaming from %q to %q: %s", from, to, err)
			}

			fd, err := Open(".")
			if err != nil {
				t.Fatalf("Open .: %s", err)
			}

			// Stat does not return the real case of the file (it returns what the called asked for)
			// So we have to use readdir to get the real name of the file.
			dirNames, err := fd.Readdirnames(-1)
			if err != nil {
				t.Fatalf("readdirnames: %s", err)
			}

			if dirNamesLen := len(dirNames); dirNamesLen != 1 {
				t.Fatalf("unexpected dirNames len, got %q, want %q", dirNamesLen, 1)
			}

			if dirNames[0] != to {
				t.Errorf("unexpected name, got %q, want %q", dirNames[0], to)
			}
		})
	}
}

func exec(t *testing.T, dir, cmd string, args []string, expect string) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe: %v", err)
	}
	defer r.Close()
	attr := &os.ProcAttr{Dir: dir, Files: []*os.File{nil, w, os.Stderr}}
	p, err := os.StartProcess(cmd, args, attr)
	if err != nil {
		t.Fatalf("StartProcess: %v", err)
	}
	w.Close()

	var b bytes.Buffer
	io.Copy(&b, r)
	output := b.String()

	fi1, _ := Stat(strings.TrimSpace(output))
	fi2, _ := Stat(expect)
	if !SameFile(fi1, fi2) {
		t.Errorf("exec %q returned %q wanted %q",
			strings.Join(append([]string{cmd}, args...), " "), output, expect)
	}
	p.Wait()
}

func TestStartProcess(t *testing.T) {
	testenv.MustHaveExec(t)

	var dir, cmd string
	var args []string
	switch runtime.GOOS {
	case "android":
		t.Skip("android doesn't have /bin/pwd")
	case "windows":
		cmd = os.Getenv("COMSPEC")
		dir = os.Getenv("SystemRoot")
		args = []string{"/c", "cd"}
	default:
		var err error
		cmd, err = osexec.LookPath("pwd")
		if err != nil {
			t.Fatalf("Can't find pwd: %v", err)
		}
		dir = "/"
		args = []string{}
		t.Logf("Testing with %v", cmd)
	}
	cmddir, cmdbase := filepath.Split(cmd)
	args = append([]string{cmdbase}, args...)
	// Test absolute executable path.
	exec(t, dir, cmd, args, dir)
	// Test relative executable path.
	exec(t, cmddir, cmdbase, args, cmddir)
}

func checkMode(t *testing.T, path string, mode FileMode) {
	dir, err := Stat(path)
	if err != nil {
		t.Fatalf("Stat %q (looking for mode %#o): %s", path, mode, err)
	}
	if dir.Mode()&0777 != mode {
		t.Errorf("Stat %q: mode %#o want %#o", path, dir.Mode(), mode)
	}
}

func TestChmod(t *testing.T) {
	// Chmod is not supported under windows.
	if runtime.GOOS == "windows" {
		return
	}
	f := newFile("TestChmod", t)
	defer Remove(f.Name())
	defer f.Close()

	if err := Chmod(f.Name(), 0456); err != nil {
		t.Fatalf("chmod %s 0456: %s", f.Name(), err)
	}
	checkMode(t, f.Name(), 0456)

	if err := f.Chmod(0123); err != nil {
		t.Fatalf("chmod %s 0123: %s", f.Name(), err)
	}
	checkMode(t, f.Name(), 0123)
}

func checkSize(t *testing.T, f Handle, size int64) {
	dir, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat %q (looking for size %d): %s", f.Name(), size, err)
	}
	if dir.Size() != size {
		t.Errorf("Stat %q: size %d want %d", f.Name(), dir.Size(), size)
	}
}

func TestFTruncate(t *testing.T) {
	f := newFile("TestFTruncate", t)
	defer Remove(f.Name())
	defer f.Close()

	checkSize(t, f, 0)
	f.Write([]byte("hello, world\n"))
	checkSize(t, f, 13)
	f.Truncate(10)
	checkSize(t, f, 10)
	f.Truncate(1024)
	checkSize(t, f, 1024)
	f.Truncate(0)
	checkSize(t, f, 0)
	_, err := f.Write([]byte("surprise!"))
	if err == nil {
		checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
	}
}

func TestTruncate(t *testing.T) {
	f := newFile("TestTruncate", t)
	defer Remove(f.Name())
	defer f.Close()

	checkSize(t, f, 0)
	f.Write([]byte("hello, world\n"))
	checkSize(t, f, 13)
	Truncate(f.Name(), 10)
	checkSize(t, f, 10)
	Truncate(f.Name(), 1024)
	checkSize(t, f, 1024)
	Truncate(f.Name(), 0)
	checkSize(t, f, 0)
	_, err := f.Write([]byte("surprise!"))
	if err == nil {
		checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
	}
}

// Use TempDir (via newFile) to make sure we're on a local file system,
// so that timings are not distorted by latency and caching.
// On NFS, timings can be off due to caching of meta-data on
// NFS servers (Issue 848).
func TestChtimes(t *testing.T) {
	f := newFile("TestChtimes", t)
	defer Remove(f.Name())

	f.Write([]byte("hello, world\n"))
	f.Close()

	testChtimes(t, f.Name())
}

// Use TempDir (via newDir) to make sure we're on a local file system,
// so that timings are not distorted by latency and caching.
// On NFS, timings can be off due to caching of meta-data on
// NFS servers (Issue 848).
func TestChtimesDir(t *testing.T) {
	name := newDir("TestChtimes", t)
	defer RemoveAll(name)

	testChtimes(t, name)
}

func testChtimes(t *testing.T, name string) {
	st, err := Stat(name)
	if err != nil {
		t.Fatalf("Stat %s: %s", name, err)
	}
	preStat := st

	// Move access and modification time back a second
	at := Atime(preStat)
	mt := preStat.ModTime()
	err = Chtimes(name, at.Add(-time.Second), mt.Add(-time.Second))
	if err != nil {
		t.Fatalf("Chtimes %s: %s", name, err)
	}

	st, err = Stat(name)
	if err != nil {
		t.Fatalf("second Stat %s: %s", name, err)
	}
	postStat := st

	pat := Atime(postStat)
	pmt := postStat.ModTime()
	if !pat.Before(at) {
		switch runtime.GOOS {
		case "plan9":
			// Mtime is the time of the last change of
			// content.  Similarly, atime is set whenever
			// the contents are accessed; also, it is set
			// whenever mtime is set.
		case "netbsd":
			mounts, _ := ioutil.ReadFile("/proc/mounts")
			if strings.Contains(string(mounts), "noatime") {
				t.Logf("AccessTime didn't go backwards, but see a filesystem mounted noatime; ignoring. Issue 19293.")
			} else {
				t.Logf("AccessTime didn't go backwards; was=%v, after=%v (Ignoring on NetBSD, assuming noatime, Issue 19293)", at, pat)
			}
		default:
			t.Errorf("AccessTime didn't go backwards; was=%v, after=%v", at, pat)
		}
	}

	if !pmt.Before(mt) {
		t.Errorf("ModTime didn't go backwards; was=%v, after=%v", mt, pmt)
	}
}

func TestFileChdir(t *testing.T) {
	// TODO(brainman): file.Chdir() is not implemented on windows.
	if runtime.GOOS == "windows" {
		return
	}

	wd, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd: %s", err)
	}
	defer Chdir(wd)

	fd, err := Open(".")
	if err != nil {
		t.Fatalf("Open .: %s", err)
	}
	defer fd.Close()

	if err := Chdir("/"); err != nil {
		t.Fatalf("Chdir /: %s", err)
	}

	if err := fd.Chdir(); err != nil {
		t.Fatalf("fd.Chdir: %s", err)
	}

	wdNew, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd: %s", err)
	}
	if wdNew != wd {
		t.Fatalf("fd.Chdir failed, got %s, want %s", wdNew, wd)
	}
}

func TestChdirAndGetwd(t *testing.T) {
	// TODO(brainman): file.Chdir() is not implemented on windows.
	/*
		align with os/os_test.go
	*/
	fd, err := Open(".")
	if err != nil {
		t.Fatalf("Open .: %s", err)
	}
	// These are chosen carefully not to be symlinks on a Mac
	// (unlike, say, /var, /etc), except /tmp, which we handle below.
	dirs := []string{"/", "/bin", "/tmp"}
	/*
		align with os/os_test.go





















	*/
	oldwd := Getenv("PWD")
	for mode := 0; mode < 2; mode++ {
		for _, d := range dirs {
			if mode == 0 {
				err = Chdir(d)
			} else {
				fd1, err1 := Open(d)
				if err1 != nil {
					t.Errorf("Open %s: %s", d, err1)
					continue
				}
				err = fd1.Chdir()
				fd1.Close()
			}
			if d == "/tmp" {
				Setenv("PWD", "/tmp")
			}
			pwd, err1 := Getwd()
			Setenv("PWD", oldwd)
			err2 := fd.Chdir()
			if err2 != nil {
				// We changed the current directory and cannot go back.
				// Don't let the tests continue; they'll scribble
				// all over some other directory.
				fmt.Fprintf(Stderr, "fchdir back to dot failed: %s\n", err2)
				Exit(1)
			}
			if err != nil {
				fd.Close()
				t.Fatalf("Chdir %s: %s", d, err)
			}
			if err1 != nil {
				fd.Close()
				t.Fatalf("Getwd in %s: %s", d, err1)
			}
			if pwd != d {
				fd.Close()
				t.Fatalf("Getwd returned %q want %q", pwd, d)
			}
		}
	}
	fd.Close()
}

// Test that Chdir+Getwd is program-wide.
func TestProgWideChdir(t *testing.T) {
	const N = 10
	const ErrPwd = "Error!"
	c := make(chan bool)
	cpwd := make(chan string, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			// Lock half the goroutines in their own operating system
			// thread to exercise more scheduler possibilities.
			if i%2 == 1 {
				// On Plan 9, after calling LockOSThread, the goroutines
				// run on different processes which don't share the working
				// directory. This used to be an issue because Go expects
				// the working directory to be program-wide.
				// See issue 9428.
				runtime.LockOSThread()
			}
			hasErr, closed := <-c
			if !closed && hasErr {
				cpwd <- ErrPwd
				return
			}
			pwd, err := Getwd()
			if err != nil {
				t.Errorf("Getwd on goroutine %d: %v", i, err)
				cpwd <- ErrPwd
				return
			}
			cpwd <- pwd
		}(i)
	}
	oldwd, err := Getwd()
	if err != nil {
		c <- true
		t.Fatalf("Getwd: %v", err)
	}
	d, err := TempDir("", "test")
	if err != nil {
		c <- true
		t.Fatalf("TempDir: %v", err)
	}
	defer func() {
		if err := Chdir(oldwd); err != nil {
			t.Fatalf("Chdir: %v", err)
		}
		RemoveAll(d)
	}()
	if err := Chdir(d); err != nil {
		c <- true
		t.Fatalf("Chdir: %v", err)
	}
	// OS X sets TMPDIR to a symbolic link.
	// So we resolve our working directory again before the test.
	d, err = Getwd()
	if err != nil {
		c <- true
		t.Fatalf("Getwd: %v", err)
	}
	close(c)
	for i := 0; i < N; i++ {
		pwd := <-cpwd
		if pwd == ErrPwd {
			t.FailNow()
		}
		if pwd != d {
			t.Errorf("Getwd returned %q; want %q", pwd, d)
		}
	}
}

func TestSeek(t *testing.T) {
	f := newFile("TestSeek", t)
	defer Remove(f.Name())
	defer f.Close()
	t.Skip("skipping test: cannot seek")
	const data = "hello, world\n"
	io.WriteString(f, data)

	type test struct {
		in     int64
		whence int
		out    int64
	}
	var tests = []test{
		{0, io.SeekCurrent, int64(len(data))},
		{0, io.SeekStart, 0},
		{5, io.SeekStart, 5},
		{0, io.SeekEnd, int64(len(data))},
		{0, io.SeekStart, 0},
		{-1, io.SeekEnd, int64(len(data)) - 1},
		{1 << 33, io.SeekStart, 1 << 33},
		{1 << 33, io.SeekEnd, 1<<33 + int64(len(data))},

		// Issue 21681, Windows 4G-1, etc:
		{1<<32 - 1, io.SeekStart, 1<<32 - 1},
		{0, io.SeekCurrent, 1<<32 - 1},
		{2<<32 - 1, io.SeekStart, 2<<32 - 1},
		{0, io.SeekCurrent, 2<<32 - 1},
	}
	for i, tt := range tests {
		off, err := f.Seek(tt.in, tt.whence)
		if off != tt.out || err != nil {
			if e, ok := err.(*PathError); ok && e != nil && tt.out > 1<<32 && runtime.GOOS == "linux" {
				mounts, _ := ioutil.ReadFile("/proc/mounts")
				if strings.Contains(string(mounts), "reiserfs") {
					// Reiserfs rejects the big seeks.
					t.Skipf("skipping test known to fail on reiserfs; https://golang.org/issue/91")
				}
			}
			t.Errorf("#%d: Seek(%v, %v) = %v, %v want %v, nil", i, tt.in, tt.whence, off, err, tt.out)
		}
	}
}

func TestSeekError(t *testing.T) {
	switch runtime.GOOS {
	case "js", "plan9":
		t.Skipf("skipping test on %v", runtime.GOOS)
	}
	t.Skip("skipping test: cannot seek")
	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Seek(0, 0)
	if err == nil {
		t.Fatal("Seek on pipe should fail")
	}
	if perr, ok := err.(*PathError); !ok || perr.Err != syscall.ESPIPE {
		t.Errorf("Seek returned error %v, want &PathError{Err: syscall.ESPIPE}", err)
	}
	_, err = w.Seek(0, 0)
	if err == nil {
		t.Fatal("Seek on pipe should fail")
	}
	if perr, ok := err.(*PathError); !ok || perr.Err != syscall.ESPIPE {
		t.Errorf("Seek returned error %v, want &PathError{Err: syscall.ESPIPE}", err)
	}
}

type openErrorTest struct {
	path  string
	mode  int
	error error
}

var openErrorTests = []openErrorTest{
	{
		sfdir + "/no-such-file",
		O_RDONLY,
		syscall.ENOENT,
	},
	{
		sfdir,
		O_WRONLY,
		syscall.EISDIR,
	},
	{
		sfdir + "/" + sfname + "/no-such-file",
		O_WRONLY,
		syscall.ENOTDIR,
	},
}

func TestOpenError(t *testing.T) {
	for _, tt := range openErrorTests {
		f, err := OpenFile(tt.path, tt.mode, 0)
		if err == nil {
			t.Errorf("Open(%q, %d) succeeded", tt.path, tt.mode)
			f.Close()
			continue
		}
		perr, ok := err.(*PathError)
		if !ok {
			t.Errorf("Open(%q, %d) returns error of %T type; want *PathError", tt.path, tt.mode, err)
		}
		if perr.Err != tt.error {
			if runtime.GOOS == "plan9" {
				syscallErrStr := perr.Err.Error()
				expectedErrStr := strings.Replace(tt.error.Error(), "file ", "", 1)
				if !strings.HasSuffix(syscallErrStr, expectedErrStr) {
					// Some Plan 9 file servers incorrectly return
					// EACCES rather than EISDIR when a directory is
					// opened for write.
					if tt.error == syscall.EISDIR && strings.HasSuffix(syscallErrStr, syscall.EACCES.Error()) {
						continue
					}
					t.Errorf("Open(%q, %d) = _, %q; want suffix %q", tt.path, tt.mode, syscallErrStr, expectedErrStr)
				}
				continue
			}
			if runtime.GOOS == "dragonfly" {
				// DragonFly incorrectly returns EACCES rather
				// EISDIR when a directory is opened for write.
				if tt.error == syscall.EISDIR && perr.Err == syscall.EACCES {
					continue
				}
			}
			t.Errorf("Open(%q, %d) = _, %q; want %q", tt.path, tt.mode, perr.Err.Error(), tt.error.Error())
		}
	}
}

func TestOpenNoName(t *testing.T) {
	f, err := Open("")
	if err == nil {
		t.Fatal(`Open("") succeeded`)
		f.Close()
	}
}

func runBinHostname(t *testing.T) string {
	t.Skip("skipping test: cannot start process")
	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	const path = "/bin/hostname"
	argv := []string{"hostname"}
	if runtime.GOOS == "aix" {
		argv = []string{"hostname", "-s"}
	}
	p, err := StartProcess(path, argv, &ProcAttr{Files: []*os.File{nil}})
	if err != nil {
		if _, err := Stat(path); IsNotExist(err) {
			t.Skipf("skipping test; test requires %s but it does not exist", path)
		}
		t.Fatal(err)
	}
	w.Close()

	var b bytes.Buffer
	io.Copy(&b, r)
	_, err = p.Wait()
	if err != nil {
		t.Fatalf("run hostname Wait: %v", err)
	}
	err = p.Kill()
	if err == nil {
		t.Errorf("expected an error from Kill running 'hostname'")
	}
	output := b.String()
	if n := len(output); n > 0 && output[n-1] == '\n' {
		output = output[0 : n-1]
	}
	if output == "" {
		t.Fatalf("/bin/hostname produced no output")
	}

	return output
}

func testWindowsHostname(t *testing.T, hostname string) {
	cmd := osexec.Command("hostname")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to execute hostname command: %v %s", err, out)
	}
	want := strings.Trim(string(out), "\r\n")
	if hostname != want {
		t.Fatalf("Hostname() = %q != system hostname of %q", hostname, want)
	}
}

func TestHostname(t *testing.T) {
	t.Skip("skipping test: cannot start process")
	hostname, err := Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if hostname == "" {
		t.Fatal("Hostname returned empty string and no error")
	}
	if strings.Contains(hostname, "\x00") {
		t.Fatalf("unexpected zero byte in hostname: %q", hostname)
	}
	// There is no other way to fetch hostname on windows, but via winapi.
	// On Plan 9 it can be taken from #c/sysname as Hostname() does.
	switch runtime.GOOS {
	case "android", "plan9":
		// No /bin/hostname to verify against.
		return
	case "windows":
		testWindowsHostname(t, hostname)
		return
	}

	testenv.MustHaveExec(t)

	// Check internal Hostname() against the output of /bin/hostname.
	// Allow that the internal Hostname returns a Fully Qualified Domain Name
	// and the /bin/hostname only returns the first component
	want := runBinHostname(t)
	if hostname != want {
		i := strings.Index(hostname, ".")
		if i < 0 || hostname[0:i] != want {
			t.Errorf("Hostname() = %q, want %q", hostname, want)
		}
	}
}

func TestReadAt(t *testing.T) {
	f := newFile("TestReadAt", t)
	defer Remove(f.Name())
	defer f.Close()

	const data = "hello, world\n"
	io.WriteString(f, data)

	b := make([]byte, 5)
	n, err := f.ReadAt(b, 7)
	if err != nil || n != len(b) {
		t.Fatalf("ReadAt 7: %d, %v", n, err)
	}
	if string(b) != "world" {
		t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
	}
}

// Verify that ReadAt doesn't affect seek offset.
// In the Plan 9 kernel, there used to be a bug in the implementation of
// the pread syscall, where the channel offset was erroneously updated after
// calling pread on a file.
func TestReadAtOffset(t *testing.T) {
	f := newFile("TestReadAtOffset", t)
	defer Remove(f.Name())
	defer f.Close()

	const data = "hello, world\n"
	io.WriteString(f, data)

	f.Seek(0, 0)
	b := make([]byte, 5)

	n, err := f.ReadAt(b, 7)
	if err != nil || n != len(b) {
		t.Fatalf("ReadAt 7: %d, %v", n, err)
	}
	if string(b) != "world" {
		t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
	}

	n, err = f.Read(b)
	if err != nil || n != len(b) {
		t.Fatalf("Read: %d, %v", n, err)
	}
	if string(b) != "hello" {
		t.Fatalf("Read: have %q want %q", string(b), "hello")
	}
}

// Verify that ReadAt doesn't allow negative offset.
func TestReadAtNegativeOffset(t *testing.T) {
	f := newFile("TestReadAtNegativeOffset", t)
	defer Remove(f.Name())
	defer f.Close()

	const data = "hello, world\n"
	io.WriteString(f, data)

	f.Seek(0, 0)
	b := make([]byte, 5)

	n, err := f.ReadAt(b, -10)

	const wantsub = "negative"
	if !strings.Contains(fmt.Sprint(err), wantsub) || n != 0 {
		t.Errorf("ReadAt(-10) = %v, %v; want 0, ...%q...", n, err, wantsub)
	}
}

func TestWriteAt(t *testing.T) {
	f := newFile("TestWriteAt", t)
	defer Remove(f.Name())
	defer f.Close()

	const data = "hello, world\n"
	io.WriteString(f, data)

	n, err := f.WriteAt([]byte("WORLD"), 7)
	if err != nil || n != 5 {
		t.Fatalf("WriteAt 7: %d, %v", n, err)
	}

	b, err := ReadFile(f.Name())
	if err != nil {
		t.Fatalf("ReadFile %s: %v", f.Name(), err)
	}
	if string(b) != "hello, WORLD\n" {
		t.Fatalf("after write: have %q want %q", string(b), "hello, WORLD\n")
	}
}

// Verify that WriteAt doesn't allow negative offset.
func TestWriteAtNegativeOffset(t *testing.T) {
	f := newFile("TestWriteAtNegativeOffset", t)
	defer Remove(f.Name())
	defer f.Close()

	n, err := f.WriteAt([]byte("WORLD"), -10)

	const wantsub = "negative offset"
	if !strings.Contains(fmt.Sprint(err), wantsub) || n != 0 {
		t.Errorf("WriteAt(-10) = %v, %v; want 0, ...%q...", n, err, wantsub)
	}
}

// Verify that WriteAt doesn't work in append mode.
func TestWriteAtInAppendMode(t *testing.T) {
	defer chtmpdir(t)()
	f, err := OpenFile("write_at_in_append_mode.txt", O_APPEND|O_CREATE, 0666)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer f.Close()

	_, err = f.WriteAt([]byte(""), 1)
	if err == nil {
		t.Fatalf("f.WriteAt returned %v, expected %v", err, ErrWriteAtInAppendMode)
	}
}

func writeFile(t *testing.T, fname string, flag int, text string) string {
	f, err := OpenFile(fname, flag, 0666)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	n, err := io.WriteString(f, text)
	if err != nil {
		t.Fatalf("WriteString: %d, %v", n, err)
	}
	f.Close()
	data, err := ReadFile(fname)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return string(data)
}

func TestAppend(t *testing.T) {
	defer chtmpdir(t)()
	const f = "append.txt"
	s := writeFile(t, f, O_CREATE|O_TRUNC|O_RDWR, "new")
	if s != "new" {
		t.Fatalf("writeFile: have %q want %q", s, "new")
	}
	s = writeFile(t, f, O_APPEND|O_RDWR, "|append")
	if s != "new|append" {
		t.Fatalf("writeFile: have %q want %q", s, "new|append")
	}
	s = writeFile(t, f, O_CREATE|O_APPEND|O_RDWR, "|append")
	if s != "new|append|append" {
		t.Fatalf("writeFile: have %q want %q", s, "new|append|append")
	}
	err := Remove(f)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	s = writeFile(t, f, O_CREATE|O_APPEND|O_RDWR, "new&append")
	if s != "new&append" {
		t.Fatalf("writeFile: after append have %q want %q", s, "new&append")
	}
	s = writeFile(t, f, O_CREATE|O_RDWR, "old")
	if s != "old&append" {
		t.Fatalf("writeFile: after create have %q want %q", s, "old&append")
	}
	s = writeFile(t, f, O_CREATE|O_TRUNC|O_RDWR, "new")
	if s != "new" {
		t.Fatalf("writeFile: after truncate have %q want %q", s, "new")
	}
}
