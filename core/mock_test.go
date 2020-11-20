package core

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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

var ErrNotExist = os.ErrNotExist

// lstat is overridden in tests.
var lstat = Lstat
var LstatP = &lstat

func Getwd() (dir string, err error) {
	return kfs.Getwd()
}

func Open(name string) (Handle, error) {
	return kfs.Open(name)
}

func Create(name string) (Handle, error) {
	return kfs.Create(name)
}

func Stat(name string) (os.FileInfo, error) {
	return kfs.Stat(name)
}

func Lstat(name string) (os.FileInfo, error) {
	return kfs.Lstat(name)
}

// Symlink creates newname as a symbolic link to oldname.
func Symlink(oldname, newname string) error {
	return e.ENotImpl
}

func RemoveAll(path string) error {
	return kfs.RemoveAll(path)
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
func Mkdir(name string, perm os.FileMode) error {
	return kfs.Mkdir(name, perm)
}

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

var errPatternHasSeparator = errors.New("pattern contains path separator")

// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
// returning prefix as the part before "*" and suffix as the part after "*".
func prefixAndSuffix(pattern string) (prefix, suffix string, err error) {
	if strings.ContainsRune(pattern, os.PathSeparator) {
		err = errPatternHasSeparator
		return
	}
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return
}

// TempDir creates a new temporary directory in the directory dir.
// The directory name is generated by taking pattern and applying a
// random string to the end. If pattern includes a "*", the random string
// replaces the last "*". TempDir returns the name of the new directory.
// If dir is the empty string, TempDir uses the
// default directory for temporary files (see os.TempDir).
// Multiple programs calling TempDir simultaneously
// will not choose the same directory. It is the caller's responsibility
// to remove the directory when no longer needed.
func TempDir(dir, pattern string) (name string, err error) {
	if dir == "" {
		dir = "/tmp"
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		try := filepath.Join(dir, prefix+nextRandom()+suffix)
		err = Mkdir(try, 0700)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		if os.IsNotExist(err) {
			if _, err := Stat(dir); os.IsNotExist(err) {
				return "", err
			}
		}
		if err == nil {
			name = try
		}
		break
	}
	return
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// The filename is generated by taking pattern and adding a random
// string to the end. If pattern includes a "*", the random string
// replaces the last "*".
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func TempFile(dir, pattern string) (f Handle, err error) {
	if dir == "" {
		dir = "/tmp"
	}

	prefix, suffix, err := prefixAndSuffix(pattern)
	if err != nil {
		return
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		f, err = kfs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

// Remove removes the named file or (empty) directory.
// If there is an error, it will be of type *PathError.
func Remove(name string) error {
	return kfs.Remove(name)
}

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func Link(oldname, newname string) error {
	return e.ENotImpl
}

// Chdir changes the current working directory to the named directory.
func Chdir(dir string) error {
	return kfs.Chdir(dir)
}

func SameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

const ModeSymlink = os.ModeSymlink

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.
func Readlink(name string) (string, error) {
	return "", e.ENotImpl
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
// If there is an error, it will be of type *LinkError.
func Rename(oldpath, newpath string) error {
	return kfs.Rename(oldpath, newpath)
}

var IsNotExist = os.IsNotExist
var IsExist = os.IsExist

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm
// (before umask); otherwise WriteFile truncates it before writing, without changing permissions.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := kfs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
