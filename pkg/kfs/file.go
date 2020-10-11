package kfs

import (
	"bytes"
	"io"

	"github.com/lazyxu/kfs/storage/obj"

	"github.com/sirupsen/logrus"
)

type File struct {
	ItemBase
}

func NewFile(kfs *KFS, name string) *File {
	return &File{
		ItemBase: ItemBase{
			kfs:      kfs,
			Metadata: obj.NewFileMetadata(name),
		},
	}
}

func (i *File) GetContent() (io.Reader, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	f := new(obj.File)
	err := f.Read(i.kfs.scheduler, i.Metadata.Hash)
	if err != nil {
		return nil, err
	}
	return f.Reader, nil
}

func (i *File) SetContent(content []byte, offset int64) (int64, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     len(content),
	}).Debug("SetContent")
	buf := make([]byte, offset)
	f := new(obj.File)
	err := f.Read(i.kfs.scheduler, i.Metadata.Hash)
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
	return i.Metadata.Size, nil
}

func (i *File) Truncate(size int64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	content := make([]byte, size)
	f := new(obj.File)
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
