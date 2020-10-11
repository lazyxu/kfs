package core

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/lazyxu/kfs/object"

	"github.com/sirupsen/logrus"
)

type File struct {
	ItemBase
}

func NewFile(kfs *KFS, name string) *File {
	return &File{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: object.NewFileMetadata(name),
		},
	}
}

func (i *File) Read(buff []byte) (int, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	reader, err := i.getContent()
	if err != nil {
		return 0, err
	}
	num, err := reader.Read(buff)
	return num, err
}

func (i *File) ReadAt(buff []byte, off int64) (int, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	reader, err := i.getContent()
	if err != nil {
		return 0, err
	}
	switch r := reader.(type) {
	case io.Seeker:
		n, err := r.Seek(off, io.SeekCurrent)
		if err != nil {
			return int(n), err
		}
	default:
		n, err := io.CopyN(ioutil.Discard, r, off)
		if err != nil {
			return int(n), err
		}
	}
	num, err := reader.Read(buff)
	return num, err
}

func (i *File) getContent() (io.Reader, error) {
	f := new(object.Blob)
	err := f.Read(i.kfs.scheduler, i.Metadata.Hash)
	if err != nil {
		return nil, err
	}
	return f.Reader, nil
}

func (i *File) WriteAt(content []byte, offset int64) (n int, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     len(content),
	}).Debug("SetContent")
	buf := make([]byte, offset)
	f := new(object.Blob)
	err = f.Read(i.kfs.scheduler, i.Metadata.Hash)
	if err != nil {
		return 0, err
	}
	if offset != 0 {
		_, err = f.Reader.Read(buf)
		if err != nil {
			return 0, err
		}
	}
	content = append(buf, content...)
	f.Reader = bytes.NewReader(content)
	hash, err := f.Write(i.kfs.scheduler)
	if err != nil {
		return 0, err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = int64(len(content))
	return int(i.Metadata.Size), nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	f := new(object.Blob)
	if size != 0 {
		err := f.Read(i.kfs.scheduler, i.Metadata.Hash)
		if err != nil {
			return err
		}
		_, err = f.Reader.Read(content)
		if err != nil {
			return err
		}
	}
	f.Reader = bytes.NewReader(content)
	hash, err := f.Write(i.kfs.scheduler)
	if err != nil {
		return err
	}
	i.Metadata.Hash = hash
	i.Metadata.Size = size
	return nil
}
