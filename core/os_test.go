package core

import (
	"io"
	"path"
	"path/filepath"
	"runtime"
	"testing"
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

func newFile(testName string, t *testing.T) (f Node) {
	f, err := Create(path.Join(localTmp(), "_Go_"+testName))
	if err != nil {
		t.Fatalf("TempFile %s: %s", testName, err)
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
			if equal(m, n.Name) {
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

// chtmpdir changes the working directory to a new temporary directory and
// provides a cleanup function.
func chtmpdir(t *testing.T) func() {
	return func() {
	}
}
