package core

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lazyxu/kfs/node"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/core/kfscommon"
)

type FS interface {
	Object() *object.Obj
	Storage() storage.Storage
	GetNode(path string) (node node.Node, err error)
}

type KFS struct {
	obj     *object.Obj
	root    *node.Dir
	storage storage.Storage
	Opt     *kfscommon.Options
	pwd     string
}

func (kfs *KFS) Object() *object.Obj {
	return kfs.obj
}
func (kfs *KFS) Storage() storage.Storage {
	return kfs.storage
}

var defaultBinaries = []string{
	"cat",
	"chmod",
	"cp",
	"date",
	"dd",
	"df",
	"hostname",
	"link",
	"ln",
	"ls",
	"mkdir",
	"mv",
	"ps",
	"pwd",
	"rm",
	"rmdir",
	"sync",
	"unlink",
}

const DevNull = "/dev/null"

func New(opt *kfscommon.Options, s storage.Storage,
	hashFunc func() kfscrypto.Hash, serializable kfscrypto.Serializable) *KFS {
	obj := object.Init(hashFunc, serializable)
	kfs := &KFS{
		Opt:     opt,
		storage: s,
		pwd:     "/tmp",
		obj:     obj,
	}
	obj.EmptyDir.Write(kfs.storage)
	obj.EmptyFile.Write(kfs.storage)
	kfs.root = node.NewDir(s, obj, obj.NewDirMetadata("", object.DefaultDirMode), nil)
	err := kfs.root.AddChild(obj.NewDirMetadata("demo", object.DefaultDirMode), obj.EmptyDir)
	if err != nil {
		panic(err)
	}
	err = kfs.root.AddChild(
		obj.NewFileMetadata("hello", object.DefaultFileMode),
		&object.Blob{Reader: strings.NewReader("hello world")})
	if err != nil {
		panic(err)
	}
	err = kfs.root.AddChild(
		obj.NewFileMetadata("index.js", object.DefaultFileMode),
		&object.Blob{Reader: strings.NewReader("index")})
	if err != nil {
		panic(err)
	}
	err = kfs.MkdirAll("/home/test", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	err = kfs.Mkdir("/bin", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	for _, b := range defaultBinaries {
		_, err = kfs.Create(path.Join("/bin", b))
		if err != nil {
			panic(err)
		}
	}
	err = kfs.Mkdir("/tmp", kfs.Opt.DirPerms)
	if err != nil {
		panic(err)
	}
	return kfs
}

func (kfs *KFS) Getwd() (dir string, err error) {
	return kfs.pwd, nil
}

func (kfs *KFS) GetDir(path string) (dir *node.Dir, err error) {
	n, err := kfs.GetNode(path)
	if err != nil {
		return nil, err
	}
	dir, ok := n.(*node.Dir)
	if !ok {
		return nil, e.ENotDir
	}
	return dir, nil
}

func (kfs *KFS) GetFile(path string) (*node.File, error) {
	n, err := kfs.GetNode(path)
	if err != nil {
		return nil, err
	}
	file, ok := n.(*node.File)
	if !ok {
		return nil, e.ENotFile
	}
	return file, nil
}

// GetNode finds the Node by path starting from the root
func (kfs *KFS) GetNode(path string) (n node.Node, err error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(kfs.pwd, path)
	}
	path = strings.Trim(path, "/")
	n = kfs.root
	for path != "" {
		i := strings.IndexRune(path, '/')
		var name string
		if i < 0 {
			name, path = path, ""
		} else {
			name, path = path[:i], path[i+1:]
		}
		if name == "" {
			continue
		}
		dir, ok := n.(*node.Dir)
		if !ok {
			// We need to look in a directory, but found a file
			return nil, e.ENotDir
		}
		n, ok = dir.Items[name]
		if ok {
			continue
		}

		d, err := kfs.obj.ReadDir(kfs.storage, dir.Metadata.Hash)
		if err != nil {
			return nil, err
		}
		metadata, err := d.GetNode(name)
		if err != nil {
			return nil, err
		}
		if metadata.IsDir() {
			n = node.NewDir(kfs.storage, kfs.obj, metadata, dir)
			dir.Items[name] = n
		} else {
			n = node.NewFile(kfs.storage, kfs.obj, metadata, dir)
			dir.Items[name] = n
		}
	}
	return
}

func (kfs *KFS) getDirAndLeaf(name string) (*node.Dir, string, error) {
	parent, leaf := filepath.Split(name)
	dir, err := kfs.GetDir(parent)
	if err != nil {
		return nil, "", err
	}
	return dir, leaf, nil
}

// Open a file according to the flags provided
//
//   O_RDONLY open the file read-only.
//   O_WRONLY open the file write-only.
//   O_RDWR   open the file read-write.
//
//   O_APPEND append data to the file when writing.
//   O_CREATE create a new file if none exists.
//   O_EXCL   used with O_CREATE, file must not exist
//   O_SYNC   open for synchronous I/O.
//   O_TRUNC  if possible, truncate file when opene
//
// We ignore O_SYNC and O_EXCL
func (kfs *KFS) _openFile(i *node.File, flags int) (fd *Handle, err error) {
	var (
		write    bool // if set need write support
		read     bool // if set need read support
		rdwrMode = flags & accessModeMask
	)

	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if rdwrMode == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		return nil, e.ErrInvalid
	}

	if flags&os.O_TRUNC != 0 {
		err := i.Truncate(0)
		if err != nil {
			return nil, err
		}
	}
	// Figure out the read/write intents
	switch {
	case rdwrMode == os.O_RDONLY:
		read = true
	case rdwrMode == os.O_WRONLY:
		write = true
	case rdwrMode == os.O_RDWR:
		read = true
		write = true
	default:
		logrus.Debug(i.Name(), "Can't figure out how to open with flags: 0x%X", flags)
		return nil, e.ErrPermission
	}

	return &Handle{
		kfs:    kfs,
		path:   i.Path(),
		read:   read,
		write:  write,
		append: flags&os.O_APPEND != 0,
	}, nil
}

// Open the directory according to the flags provided
func (kfs *KFS) _openDir(d *node.Dir, flags int) (fd *Handle, err error) {
	rdwrMode := flags & accessModeMask
	if rdwrMode != os.O_RDONLY {
		logrus.Error(d, "Can only open directories read only")
		return nil, e.EIsDir
	}
	return &Handle{
		kfs:   kfs,
		path:  d.Path(),
		isDir: true,
		read:  true,
	}, nil
}
