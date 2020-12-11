package core

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lazyxu/kfs/node"

	"github.com/lazyxu/kfs/core/e"
)

var ErrNotExist = os.ErrNotExist
var ErrClosed = os.ErrClosed

// lstat is overridden in tests.
var lstat = Lstat
var LstatP = &lstat

func Getwd() (dir string, err error) {
	return kfs.Getwd()
}

func Open(name string) (*Handle, error) {
	return kfs.Open(name)
}

func Create(name string) (*Handle, error) {
	return kfs.Create(name)
}

func Stat(name string) (os.FileInfo, error) {
	info, err := kfs.Stat(name)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}
	return info, nil
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
func TempFile(dir, pattern string) (f *Handle, err error) {
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
	f1, ok1 := fi1.(*node.File)
	f2, ok2 := fi2.(*node.File)
	if !ok1 || !ok2 {
		return false
	}
	return f1.Metadata() == f2.Metadata()
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
	err := kfs.Rename(oldpath, newpath)
	if err != nil {
		return &LinkError{"rename", oldpath, newpath, err}
	}
	return nil
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

type FileMode = os.FileMode
type FileInfo = os.FileInfo

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.
//
// A different subset of the mode bits are used, depending on the
// operating system.
//
// On Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and
// ModeSticky are used.
//
// On Windows, only the 0200 bit (owner writable) of mode is used; it
// controls whether the file's read-only attribute is set or cleared.
// The other bits are currently unused. For compatibility with Go 1.12
// and earlier, use a non-zero mode. Use mode 0400 for a read-only
// file and 0600 for a readable+writable file.
//
// On Plan 9, the mode's permission bits, ModeAppend, ModeExclusive,
// and ModeTemporary are used.
func Chmod(name string, mode FileMode) error { return kfs.Chmod(name, mode) }

// Truncate changes the size of the named file.
// If the file is a symbolic link, it changes the size of the link's target.
// If there is an error, it will be of type *PathError.
func Truncate(name string, size int64) error {
	if e := kfs.Truncate(name, size); e != nil {
		return &PathError{"truncate", name, e}
	}
	return nil
}
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// For testing.
func Atime(fi os.FileInfo) time.Time {
	n, ok := fi.(node.Node)
	if !ok {
		panic(errors.New("235"))
	}
	return n.ModTime()
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func Chtimes(name string, atime time.Time, mtime time.Time) error {
	node, err := kfs.GetNode(name)
	if err != nil {
		return err
	}
	node.SetATime(atime)
	node.SetMTime(mtime)
	return nil
}

var Setenv = os.Setenv
var Getenv = os.Getenv
var Stderr = os.Stderr
var Exit = os.Exit

// Pipe returns a connected pair of Files; reads from r return bytes written to w.
// It returns the files and an error, if any.
func Pipe() (r *Handle, w *Handle, err error) {
	return nil, nil, e.ENotImpl
}

const O_RDONLY = os.O_RDONLY
const O_WRONLY = os.O_WRONLY

// OpenFile is the generalized open call; most users will use Open
// or Create instead. It opens the named file with specified flag
// (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
// is passed, it is created with mode perm (before umask). If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFile(name string, flag int, perm FileMode) (*Handle, error) {
	h, err := kfs.OpenFile(name, flag, perm)
	if err != nil {
		return nil, wrapErr("open", name, err)
	}
	return h, nil
}

var StartProcess = os.StartProcess

type ProcAttr = os.ProcAttr

var Hostname = os.Hostname

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func ReadFile(filename string) ([]byte, error) {
	n, err := kfs.GetFile(filename)
	if err != nil {
		return nil, err
	}
	return n.ReadAll()
}

const O_APPEND = os.O_APPEND
const O_CREATE = os.O_CREATE
const O_TRUNC = os.O_TRUNC
const O_RDWR = os.O_RDWR

type ProcessState = os.ProcessState

const ModeDevice = os.ModeDevice
const ModeCharDevice = os.ModeCharDevice

var Stdout = os.Stdout
var Stdin = os.Stdin
var ModeNamedPipe = os.ModeNamedPipe
var Environ = os.Environ
var Args = os.Args

func MkdirAll(path string, perm os.FileMode) error {
	return kfs.MkdirAll(path, perm)
}

type Process = os.Process

var Getppid = os.Getppid
var Getpid = os.Getpid
var FindProcess = os.FindProcess
var ErrInvalid = os.ErrInvalid
var IsTimeout = os.IsTimeout
var ErrDeadlineExceeded = os.ErrDeadlineExceeded

func UserHomeDir() (string, error) {
	return "/home", nil
}
