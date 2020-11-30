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
			Metadata: metadata,
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
	reader, err := i.getContent()
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
	reader, err := i.getContent()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (i *File) getContent() (io.Reader, error) {
	blob := new(object.Blob)
	err := blob.Read(i.storage, i.Metadata.Hash)
	if err != nil {
		return nil, err
	}
	return blob.Reader, nil
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
	blob := new(object.Blob)
	err = blob.Read(i.storage, i.Metadata.Hash)
	if err != nil {
		return 0, err
	}
	if offset != 0 {
		_, err = blob.Reader.Read(buf)
		if err != nil {
			return 0, err
		}
	}
	content = append(buf, content...)
	n, err = skip(blob.Reader, int64(l))
	if err != nil && err != io.EOF {
		return n, err
	}
	if err != io.EOF {
		remain, err := ioutil.ReadAll(blob.Reader)
		if err != nil {
			return 0, err
		}
		content = append(content, remain...)
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.storage)
	if err != nil {
		return 0, err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = int64(len(content))
	return l, nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	blob := new(object.Blob)
	if size != 0 {
		err := blob.Read(i.storage, i.Metadata.Hash)
		if err != nil {
			return err
		}
		_, err = blob.Reader.Read(content)
		if err != nil {
			return err
		}
	}
	blob.Reader = bytes.NewReader(content)
	hash, err := blob.Write(i.storage)
	if err != nil {
		return err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = size
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
