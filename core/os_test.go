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
	"testing"
	"time"
)

var dot = []string{
	"dir.go",
	"file.go",
	"node.go",
	"os_test.go",
	"stat.go",
}

type sysDir struct {
	name  string
	files []string
}

var sysdir = &sysDir{
	"/etc",
	[]string{
		"group",
		"hosts",
		"passwd",
	},
}

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
	return name1 == name2
}

// localTmp returns a local temporary directory not on NFS.
func localTmp() string {
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
//func TestReadClosed(t *testing.T) {
//	path := sfdir + "/" + sfname
//	file, err := Open(path)
//	if err != nil {
//		t.Fatal("open failed:", err)
//	}
//	file.Close() // close immediately
//
//	b := make([]byte, 100)
//	_, err = file.Read(b)
//
//	if err != os.ErrClosed {
//		t.Errorf("Read: %v, want ErrClosed", err)
//	}
//}

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
	if !IsExist(err) {
		t.Errorf("Link(%q, %q) err.Err = %q; want %q", to, from, err, "file exists error")
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
	if err == nil {
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
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
	if err == nil {
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
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
