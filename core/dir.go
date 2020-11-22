package core

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/e"
)

type Dir struct {
	ItemBase
	items map[string]Node
}

func NewDir(kfs *KFS, name string, perm os.FileMode) *Dir {
	return &Dir{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: object.NewDirMetadata(name, perm),
		},
		items: make(map[string]Node),
	}
}

func (i *Dir) load() (*object.Tree, error) {
	tree := new(object.Tree)
	err := tree.Read(i.kfs.scheduler, i.Metadata.Hash)
	return tree, err
}

func getSize(r io.Reader) (int64, error) {
	switch v := r.(type) {
	case *bytes.Buffer:
		return int64(v.Len()), nil
	case *bytes.Reader:
		return int64(v.Len()), nil
	case *strings.Reader:
		return int64(v.Len()), nil
	case *os.File:
		info, err := v.Stat()
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	default:
		return 0, errors.New("invalid type: " + reflect.TypeOf(r).String())
	}
}

func (i *Dir) add(metadata *object.Metadata, item object.Object) error {
	d, err := i.load()
	if err != nil {
		return err
	}
	for _, it := range d.Items {
		if it.Name == metadata.Name {
			return e.ErrExist
		}
	}

	if blob, ok := item.(*object.Blob); ok {
		size, err := getSize(blob.Reader)
		if err != nil {
			return err
		}
		metadata.Size = size
	}
	itemHash, err := item.Write(i.kfs.scheduler)
	if err != nil {
		return err
	}
	metadata.Hash = itemHash
	d.Items = append(d.Items, metadata)

	return i.updateObj(d)
}

func (i *Dir) Create(name string, flags int) (*File, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	f := &File{
		ItemBase: ItemBase{
			kfs:      i.kfs,
			parent:   i,
			Metadata: object.NewFileMetadata(name),
		},
	}
	err := i.add(f.Metadata, object.EmptyFile)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (i *Dir) get(name string) (*object.Metadata, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	for _, item := range d.Items {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, e.ErrNotExist
}

func (i *Dir) remove(name string, all bool) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	d, err := i.load()
	if err != nil {
		return err
	}
	for index, item := range d.Items {
		if item.Name == name {
			if all || item.IsFile() || item.Hash == object.EmptyDirHash {
				d.Items = append(d.Items[0:index], d.Items[index+1:]...)
				delete(i.items, name)
				return i.updateObj(d)
			}
			if item.IsDir() {
				return e.ENotEmpty
			}
		}
	}
	if all {
		return nil
	}
	return e.ErrNotExist
}

func (i *Dir) ReadDirAll() ([]*object.Metadata, error) {
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	return d.Items, nil
}

func (i *Dir) Read(buff []byte) (int, error) {
	return 0, &os.PathError{
		Op:   "read",
		Path: i.Name(),
		Err:  e.EIsDir,
	}
}

func (i *Dir) ReadAt(buff []byte, off int64) (int, error) {
	return 0, &os.PathError{
		Op:   "read",
		Path: i.Name(),
		Err:  e.EIsDir,
	}
}

func (i *Dir) Write(content []byte) (n int, err error) {
	return i.WriteAt(content, 0)
}

func (i *Dir) WriteAt(content []byte, offset int64) (n int, err error) {
	return 0, &os.PathError{
		Op:   "write",
		Path: i.Name(),
		Err:  e.EIsDir,
	}
}

// Readdir reads the contents of the directory associated with file and
// returns a slice of up to n FileInfo values, as would be returned
// by Lstat, in directory order. Subsequent calls on the same file will yield
// further FileInfos.
//
// If n > 0, Readdir returns at most n FileInfo structures. In this case, if
// Readdir returns an empty slice, it will return a non-nil error
// explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdir returns all the FileInfo from the directory in
// a single slice. In this case, if Readdir succeeds (reads all
// the way to the end of the directory), it returns the slice and a
// nil error. If it encounters an error before the end of the
// directory, Readdir returns the FileInfo read until that point
// and a non-nil error.
func (i *Dir) Readdir(n int, offset int) (dirs []*object.Metadata, err error) {
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		if offset >= len(d.Items) {
			return []*object.Metadata{}, nil
		}
		offset = len(d.Items)
		return d.Items, nil
	}
	if offset >= len(d.Items) {
		return []*object.Metadata{}, io.EOF
	}
	for ii := offset; ii < len(d.Items); ii++ {
		if ii >= offset+n {
			break
		}
		dirs = append(dirs, d.Items[ii])
	}
	return dirs, nil
}

// Readdirnames reads the contents of the directory associated with file
// and returns a slice of up to n names of files in the directory,
// in directory order. Subsequent calls on the same file will yield
// further names.
//
// If n > 0, Readdirnames returns at most n names. In this case, if
// Readdirnames returns an empty slice, it will return a non-nil error
// explaining why. At the end of a directory, the error is io.EOF.
//
// If n <= 0, Readdirnames returns all the names from the directory in
// a single slice. In this case, if Readdirnames succeeds (reads all
// the way to the end of the directory), it returns the slice and a
// nil error. If it encounters an error before the end of the
// directory, Readdirnames returns the names read until that point and
// a non-nil error.
func (i *Dir) Readdirnames(n int, nameOffset int) (names []string, err error) {
	d, err := i.load()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		if nameOffset >= len(d.Items) {
			return []string{}, nil
		}
		names = make([]string, len(d.Items))
		for ii, item := range d.Items {
			names[ii] = item.Name
		}
		nameOffset = len(d.Items)
		return names, nil
	}
	if nameOffset >= len(d.Items) {
		return []string{}, io.EOF
	}
	for ii := nameOffset; ii < len(d.Items); ii++ {
		if ii >= nameOffset+n {
			break
		}
		names = append(names, d.Items[ii].Name)
	}
	return names, nil
}

func (i *Dir) Close() error {
	err := i.ItemBase.Close()
	if err != nil {
		return err
	}
	return nil
}

// Open the directory according to the flags provided
func (d *Dir) Open(flags int) (fd Handle, err error) {
	//rdwrMode := flags & accessModeMask
	//if rdwrMode != os.O_RDONLY {
	//	logrus.Error(d, "Can only open directories read only")
	//	return nil, e.ErrPermission
	//}
	return newDirHandle(d.kfs, d.Path()), nil
}

func (i *Dir) Truncate(size int64) error {
	return e.ENotImpl
}
