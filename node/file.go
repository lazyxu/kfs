package node

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/core/e"

	"github.com/lazyxu/kfs/object"

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
func skip(reader io.Reader, off int64) (int, error) {
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		return int(n), err
	}
	n, err := io.CopyN(ioutil.Discard, reader, off)
	return int(n), err
}

func (i *File) ReadAt(buff []byte, off int64) (int, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	if len(buff) == 0 {
		return 0, nil
	}
	reader, err := i.Content()
	if err != nil {
		return 0, err
	}
	n, err := skip(reader, off)
	if err != nil {
		return n, err
	}
	num, err := reader.Read(buff)
	return num, err
}

func (i *File) ReadAll() ([]byte, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	reader, err := i.Content()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (i *File) Content() (io.Reader, error) {
	r, err := i.obj.ReadBlob(i.metadata.Hash)
	if err != nil {
		return nil, err
	}
	return r, nil
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
	r, err := i.obj.ReadBlob(i.metadata.Hash)
	if err != nil {
		return 0, err
	}
	if offset != 0 {
		_, err = r.Read(buf)
		if err != nil {
			return 0, err
		}
	}
	content = append(buf, content...)
	n, err = skip(r, int64(l))
	if err != nil && err != io.EOF {
		return n, err
	}
	if err != io.EOF {
		remain, err := ioutil.ReadAll(r)
		if err != nil {
			return 0, err
		}
		content = append(content, remain...)
	}
	hash, err := i.obj.WriteBlob(bytes.NewReader(content))
	if err != nil {
		return 0, err
	}
	i.metadata.Hash = hash
	i.metadata.Size = int64(len(content))
	return l, nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	blob := i.obj.NewBlob()
	if size != 0 {
		err := blob.Read(i.metadata.Hash)
		if err != nil {
			return err
		}
		_, err = blob.Reader.Read(content)
		if err != nil {
			return err
		}
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write()
	if err != nil {
		return err
	}
	i.metadata.Hash = hash
	i.metadata.Size = size
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
