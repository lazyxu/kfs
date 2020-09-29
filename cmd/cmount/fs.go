package main

import (
	"fmt"
	"os"

	"github.com/lazyxu/kfs/kfs/e"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/kfs/kfscommon"

	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/lazyxu/kfs/kfs"
	"github.com/lazyxu/kfs/node"
)

const (
	filename = "hello"
	contents = "hello, world\n"
)

type FS struct {
	fuse.FileSystemBase
	kfs *kfs.KFS
}

func NewFS() *FS {
	logrus.SetLevel(logrus.TraceLevel)
	return &FS{
		kfs: kfs.New(&kfscommon.Options{
			UID:       uint32(os.Getuid()),
			GID:       uint32(os.Getgid()),
			DirPerms:  fuse.S_IFDIR | 0755,
			FilePerms: fuse.S_IFREG | 0644,
		}),
	}
}

// Statfs gets file system statistics.
// The FileSystemBase implementation returns -ENotImpl.
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

func (fs *FS) Unlink(filepath string) (errCodeode int) {
	defer e.Trace(logrus.Fields{
		"path": filepath,
	})(func() logrus.Fields {
		return logrus.Fields{
			"errCodeode": errCodeode,
		}
	})
	err := fs.kfs.Remove(filepath)
	return translateError(err)
}

func (fs *FS) Open(path string, flags int) (errCodeode int, fh uint64) {
	defer e.Trace(logrus.Fields{
		"path":  path,
		"flags": flags,
	})(func() logrus.Fields {
		return logrus.Fields{
			"errCodeode": errCodeode,
			"fh":         fh,
		}
	})
	_, err := fs.kfs.OpenFile(path, flags, fs.kfs.Opt.FilePerms)
	return translateError(err), defaultFileHandler
}

func (fs *FS) Access(path string, mask uint32) int {
	defer e.Trace(logrus.Fields{
		"path": path,
		"mask": mask,
	})(nil)
	return -fuse.ENOSYS
}

func (fs *FS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errCodeode int) {
	defer e.Trace(logrus.Fields{"path": path, "fh": fh})(func() logrus.Fields {
		return logrus.Fields{
			"stat":       fmt.Sprintf("%+v", stat),
			"errCodeode": errCodeode,
		}
	})
	n, err := fs.kfs.GetNode(path)
	if err != nil {
		return translateError(err)
	}
	fs.stat(n, stat)
	return
}

func (fs *FS) Read(path string, buff []byte, off int64, fh uint64) (num int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"num": num,
		}
	})
	n, err := fs.kfs.Read(path, buff, off)
	if err != nil {
		return translateError(err)
	}
	return int(n)
}

func (fs *FS) Create(filepath string, flags int, mode uint32) (errCodeode int, fh uint64) {
	defer e.Trace(logrus.Fields{"path": filepath})(func() logrus.Fields {
		return logrus.Fields{
			"errCodeode": errCodeode,
		}
	})
	_, err := fs.kfs.OpenFile(filepath, flags, os.FileMode(mode))
	errCodeode = translateError(err)
	return
}

func (fs *FS) Write(path string, buff []byte, offset int64, fh uint64) (errCodeode int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"errCodeode": errCodeode,
		}
	})
	n, err := fs.kfs.Write(path, buff, offset)
	if err != nil {
		return translateError(err)
	}
	return int(n)
}

func (fs *FS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, offset int64) bool,
	offset int64, fh uint64) (errCodeode int) {
	defer e.Trace(logrus.Fields{"path": path})(func() logrus.Fields {
		return logrus.Fields{
			"errCodeode": errCodeode,
		}
	})
	fill(".", nil, 0)
	fill("..", nil, 0)
	nodes, err := fs.kfs.Readdir(path)
	if err != nil {
		return translateError(err)
	}
	for _, n := range nodes {
		logrus.WithFields(logrus.Fields{
			"name": n.Name(),
		}).Debug("node")
		var stat fuse.Stat_t
		fs.stat(n, &stat)
		fill(n.Name(), &stat, 0)
	}
	return 0
}

// stat fills up the stat block for Node
func (fs *FS) stat(node node.Node, stat *fuse.Stat_t) {
	size, err := node.Size()
	if err != nil {
		logrus.WithError(err).Error("node.Size()")
	}
	blocks := (size + 511) / 512
	// stat.Dev // Device ID of device containing file. [IGNORED]
	// stat.Ino // File serial number. [IGNORED unless the use_ino mount option is given.]
	stat.Mode = uint32(node.Mode())
	stat.Nlink = 1
	stat.Uid = fs.kfs.Opt.UID
	stat.Gid = fs.kfs.Opt.GID
	// stat.Rdev // Device ID (if file is character or block special).
	stat.Size = size
	stat.Atim = fuse.NewTimespec(node.AccessTime())
	stat.Mtim = fuse.NewTimespec(node.AccessTime())
	stat.Ctim = fuse.NewTimespec(node.AccessTime())
	stat.Blksize = 512
	stat.Blocks = int64(blocks)
	stat.Birthtim = fuse.NewTimespec(node.AccessTime())
	// stat.Flags
}

// Truncate truncates a file to size
func (fs *FS) Truncate(path string, size int64, fh uint64) (errCode int) {
	defer e.Trace(logrus.Fields{"path": path, "size": size, "fh": fh})(func() logrus.Fields {
		return logrus.Fields{
			"errCode": errCode,
		}
	})
	n, err := fs.kfs.GetNode(path)
	if err != nil {
		return translateError(err)
	}
	err = n.Truncate(uint64(size))
	if err != nil {
		return translateError(err)
	}
	return 0
}
