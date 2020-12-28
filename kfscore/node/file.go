package node

import (
	"bytes"
	"io"

	"github.com/lazyxu/kfs/kfscore/util"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/lazyxu/kfs/kfscore/e"

	"github.com/lazyxu/kfs/kfscore/object"

	"github.com/sirupsen/logrus"
)

type File struct {
	ItemBase
}

func NewFile(s storage.Storage, obj *object.Obj, metadata *object.Metadata, parent *Dir) *File {
	return &File{
		ItemBase: ItemBase{
			storage:  s,
			obj:      obj,
			metadata: metadata,
			Parent:   parent,
		},
	}
}

func (i *File) ReadAt(buff []byte, off int64) (int, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	if len(buff) == 0 {
		return 0, nil
	}
	var n int
	err := i.Content(func(r io.Reader) error {
		_, err := util.Skip(r, off)
		if err != nil {
			return err
		}
		n, err = r.Read(buff)
		return err
	})
	return n, err
}

func (i *File) ReadAll() ([]byte, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	buf := new(bytes.Buffer)
	err := i.Content(func(r io.Reader) error {
		_, err := io.Copy(buf, r)
		return err
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *File) Content(f func(reader io.Reader) error) error {
	return i.obj.ReadBlob(i.metadata.Hash(), f)
}

func (i *File) WriteAt(content []byte, offset int64) (n int, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	l := len(content)
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     l,
	}).Debug("SetContent")
	if offset < 0 {
		return 0, e.ENegative
	}
	buf := make([]byte, offset)
	err = i.obj.ReadBlob(i.metadata.Hash(), func(r io.Reader) error {
		if offset != 0 {
			n, err = r.Read(buf)
			if err != nil {
				return err
			}
			buf = buf[:n]
		}
		n, err = util.Skip(r, int64(l))
		if err != nil && err != io.EOF {
			return err
		}
		hash, err := i.obj.WriteBlob(io.MultiReader(
			bytes.NewReader(buf),
			bytes.NewReader(content),
			r))
		if err != nil {
			return err
		}
		remain, err := util.Size(r)
		if err != nil {
			return err
		}
		i.metadata = i.metadata.Builder().
			Hash(hash).Size(int64(len(buf)+len(content)) + remain).Build()
		return nil
	})
	return l, err
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	if size != 0 {
		err := i.obj.ReadBlob(i.metadata.Hash(), func(r io.Reader) error {
			_, err := r.Read(content)
			return err
		})
		if err != nil {
			return err
		}
	}
	hash, err := i.obj.WriteBlob(bytes.NewReader(content))
	if err != nil {
		return err
	}
	i.metadata = i.metadata.Builder().
		Hash(hash).Size(size).Build()
	return nil
}

func (i *File) Readdirnames(n int, offset int) (names []string, err error) {
	if i == nil {
		return nil, e.ErrInvalid
	}
	return nil, e.EIsFile
}

func (i *File) Readdir(n int, offset int) ([]*object.Metadata, error) {
	if i == nil {
		return nil, e.ErrInvalid
	}
	return nil, e.EIsFile
}

func (i *File) Close() error {
	err := i.ItemBase.Close()
	if err != nil {
		return err
	}
	return nil
}
