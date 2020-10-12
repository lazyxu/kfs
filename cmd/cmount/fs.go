package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/core/kfscommon"

	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/lazyxu/kfs/core"
)

type FS struct {
	kfs *core.KFS
}

func NewFS() *FS {
	logrus.SetLevel(logrus.TraceLevel)
	return &FS{
		kfs: core.New(&kfscommon.Options{
			UID:       uint32(os.Getuid()),
			GID:       uint32(os.Getgid()),
			DirPerms:  fuse.S_IFDIR | 0755,
			FilePerms: fuse.S_IFREG | 0644,
		}),
	}
}

// Init is called when the file system is created.
func (fs *FS) Init() {
}

// Destroy is called when the file system is destroyed.
func (fs *FS) Destroy() {
}

// Statfs gets file system statistics.
func (fs *FS) Statfs(path string, stat *fuse.Statfs_t) int {
	defer e.Trace(logrus.Fields{"path": path})(nil)
	const blockSize = 4096
	total := 0
	free := 0
	//total, _, free := fsys.VFS.Statfs()
	stat.Blocks = uint64(total) / blockSize // Total data blocks in file system.
	stat.Bfree = uint64(free) / blockSize   // Free blocks in file system.
	stat.Bavail = stat.Bfree                // Free blocks in file system if you're not root.
	stat.Files = 1e9                        // Total files in file system.
	stat.Ffree = 1e9                        // Free files in file system.
	stat.Bsize = blockSize                  // Block size
	stat.Namemax = 255                      // Maximum file name length?
	stat.Frsize = blockSize                 // Fragment size, smallest addressable data size in the file system.
	//mountlib.ClipBlocks(&stat.Blocks)
	//mountlib.ClipBlocks(&stat.Bfree)
	//mountlib.ClipBlocks(&stat.Bavail)
	return 0
}

// Mknod creates a file node.
func (fs *FS) Mknod(path string, mode uint32, dev uint64) int {
	return translateError(e.ENotImpl)
}

// Mkdir creates a directory.
func (fs *FS) Mkdir(path string, mode uint32) int {
	err := fs.kfs.Mkdir(path, os.FileMode(mode))
	return translateError(err)
}

// Unlink removes a file.
func (fs *FS) Unlink(filepath string) (errCode int) {
	err := fs.kfs.Remove(filepath)
	return translateError(err)
}

// Rmdir removes a directory.
func (fs *FS) Rmdir(path string) int {
	err := fs.kfs.Remove(path)
	return translateError(err)
}

// Link creates a hard link to a file.
func (fs *FS) Link(oldpath string, newpath string) int {
	return translateError(e.ENotImpl)
}

// Symlink creates a symbolic link.
func (fs *FS) Symlink(target string, newpath string) int {
	return translateError(e.ENotImpl)
}

// Readlink reads the target of a symbolic link.
func (fs *FS) Readlink(path string) (int, string) {
	return translateError(e.ENotImpl), ""
}

// Rename renames a file.
func (fs *FS) Rename(oldpath string, newpath string) int {
	err := fs.kfs.Rename(oldpath, newpath)
	return translateError(err)
}

// Chmod changes the permission bits of a file.
func (fs *FS) Chown(path string, uid uint32, gid uint32) int {
	return translateError(e.ENotImpl)
}

// Chmod changes the permission bits of a file.
func (fs *FS) Chmod(path string, mode uint32) int {
	err := fs.kfs.Chmod(path, os.FileMode(mode))
	return translateError(err)
}

// Utimens changes the access and modification times of a file.
func (fs *FS) Utimens(path string, tmsp []fuse.Timespec) int {
	return translateError(e.ENotImpl)
}

// Access checks file access permissions.
func (fs *FS) Access(path string, mask uint32) int {
	defer e.Trace(logrus.Fields{
		"path": path,
		"mask": mask,
	})(nil)
	return -fuse.ENOSYS
}

// Create creates and opens a file.
// The flags are a combination of the fuse.O_* constants.
func (fs *FS) Create(filepath string, flags int, mode uint32) (errCode int, fh uint64) {
	defer e.Trace(logrus.Fields{"path": filepath})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
		}
	})
	_, err := fs.kfs.OpenFile(filepath, flags, os.FileMode(mode))
	errCode = translateError(err)
	return
}

// Open opens a file.
// The flags are a combination of the fuse.O_* constants.
func (fs *FS) Open(path string, flags int) (errCode int, fh uint64) {
	defer e.Trace(logrus.Fields{
		"path":  path,
		"flags": flags,
	})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
			"fh":      fh,
		}
	})
	_, err := fs.kfs.OpenFile(path, flags, fs.kfs.Opt.FilePerms)
	return translateError(err), defaultFileHandler
}

// Getattr gets file attributes.
func (fs *FS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errCode int) {
	defer e.Trace(logrus.Fields{"path": path, "fh": fh})(func() logrus.Fields {
		return logrus.Fields{
			"stat":    fmt.Sprintf("%+v", stat),
			"errCode": errCode,
		}
	})
	n, err := fs.kfs.Stat(path)
	if err != nil {
		return translateError(err)
	}
	fs.stat(n.Sys().(*object.Metadata), stat)
	return
}

// Truncate changes the size of a file.
func (fs *FS) Truncate(path string, size int64, fh uint64) (errCode int) {
	defer e.Trace(logrus.Fields{"path": path, "size": size, "fh": fh})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
		}
	})
	n, err := fs.kfs.GetFile(path)
	if err != nil {
		return translateError(err)
	}
	err = n.Truncate(size)
	if err != nil {
		return translateError(err)
	}
	return 0
}

// Read reads data from a file.
func (fs *FS) Read(path string, buff []byte, off int64, fh uint64) (num int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"num": num,
		}
	})
	n, err := fs.kfs.Open(path)
	if err != nil {
		return translateError(err)
	}
	num, err = n.ReadAt(buff, off)
	if err != nil {
		return translateError(err)
	}
	return num
}

// Write writes data to a file.
func (fs *FS) Write(path string, buff []byte, offset int64, fh uint64) (errCode int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
		}
	})
	n, err := fs.kfs.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return translateError(err)
	}
	num, err := n.WriteAt(buff, offset)
	if err != nil {
		return translateError(err)
	}
	return num
}

// Flush flushes cached file data.
func (fs *FS) Flush(path string, fh uint64) int {
	return translateError(e.ENotImpl)
}

// Release closes an open file.
func (fs *FS) Release(path string, fh uint64) int {
	return translateError(e.ENotImpl)
}

// Fsync synchronizes file contents.
func (fs *FS) Fsync(path string, datasync bool, fh uint64) int {
	return translateError(e.ENotImpl)
}

// Opendir opens a directory.
func (fs *FS) Opendir(path string) (int, uint64) {
	return translateError(e.ENotImpl), ^uint64(0)
}

// Readdir reads a directory.
func (fs *FS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, offset int64) bool,
	offset int64, fh uint64) (errCode int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
		}
	})
	fill(".", nil, 0)
	fill("..", nil, 0)
	nodes, err := fs.kfs.ReadDir(path)
	if err != nil {
		return translateError(err)
	}
	for _, n := range nodes {
		logrus.WithFields(logrus.Fields{
			"name": n.Name,
		}).Debug("node")
		var stat fuse.Stat_t
		fs.stat(n, &stat)
		fill(n.Name, &stat, 0)
	}
	return 0
}

// Releasedir closes an open directory.
func (fs *FS) Releasedir(path string, fh uint64) int {
	return translateError(e.ENotImpl)
}

// Fsyncdir synchronizes directory contents.
func (fs *FS) Fsyncdir(path string, datasync bool, fh uint64) int {
	return translateError(e.ENotImpl)
}

// Setxattr sets extended attributes.
func (fs *FS) Setxattr(path string, name string, value []byte, flags int) int {
	return translateError(e.ENotImpl)
}

// Getxattr gets extended attributes.
func (fs *FS) Getxattr(path string, name string) (int, []byte) {
	return translateError(e.ENotImpl), nil
}

// Removexattr removes extended attributes.
func (fs *FS) Removexattr(path string, name string) int {
	return translateError(e.ENotImpl)
}

// Listxattr lists extended attributes.
func (fs *FS) Listxattr(path string, fill func(name string) bool) int {
	return translateError(e.ENotImpl)
}

// stat fills up the stat block for Node.
func (fs *FS) stat(metadata *object.Metadata, stat *fuse.Stat_t) {
	size := metadata.Size
	blocks := (size + 511) / 512
	// stat.Dev // Device ID of device containing file. [IGNORED]
	// stat.Ino // File serial number. [IGNORED unless the use_ino mount option is given.]
	stat.Mode = uint32(metadata.Mode)
	stat.Nlink = 1
	stat.Uid = fs.kfs.Opt.UID
	stat.Gid = fs.kfs.Opt.GID
	// stat.Rdev // Device ID (if file is character or block special).
	stat.Size = size
	stat.Atim = fuse.NewTimespec(time.Unix(0, metadata.ModifyTime))
	stat.Mtim = fuse.NewTimespec(time.Unix(0, metadata.ModifyTime))
	stat.Ctim = fuse.NewTimespec(time.Unix(0, metadata.ChangeTime))
	stat.Blksize = 512
	stat.Blocks = blocks
	stat.Birthtim = fuse.NewTimespec(time.Unix(0, metadata.BirthTime))
	// stat.Flags
}
