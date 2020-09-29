package kfs

import (
	"os"
	"sync"
	"time"

	"github.com/lazyxu/kfs/storage/scheduler"

	"github.com/lazyxu/kfs/object"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/node"
)

type File struct {
	node.TimeImpl
	kfs   *KFS
	name  string
	path  string
	hash  string
	mutex sync.RWMutex // protects the following
}

func NewFile(kfs *KFS, name string) (*File, error) {
	now := time.Now()
	file := &File{
		TimeImpl: node.TimeImpl{
			BTime: now,
			ATime: now,
			Mtime: now,
			CTime: now,
		},
		kfs:  kfs,
		name: name,
		path: name,
		hash: scheduler.EmptyFileHash,
	}
	return file, nil
}

func (f *File) Name() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.name
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) IsFile() bool {
	return true
}

func (f *File) GetContent() (string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	fileObj, err := f.kfs.scheduler.GetFileObjectByHash(f.Hash())
	if err != nil {
		return "", err
	}
	return fileObj.Content, nil
}

func (f *File) SetContent(content []byte, offset int64) (int64, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     len(content),
	}).Debug("SetContent")
	fileObj, err := f.kfs.scheduler.GetFileObjectByHash(f.Hash())
	if err != nil {
		return 0, err
	}
	content = append([]byte(fileObj.Content)[0:offset], content...)
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"len":     len(content),
		"path":    f.path,
	}).Debug("AfterSetContent")
	newFile := &object.File{Content: string(content)}
	err = f.kfs.scheduler.SetObjectByHash(newFile)
	if err != nil {
		return 0, err
	}
	f.hash = newFile.Hash()
	return f.Size()
}

func (f *File) Size() (int64, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	fileObj, err := f.kfs.scheduler.GetFileObjectByHash(f.Hash())
	if err != nil {
		return 0, err
	}
	return int64(len(fileObj.Content)), nil
}

func (f *File) Mode() (mode os.FileMode) {
	return f.kfs.Opt.FilePerms
}

func (f *File) Truncate(size uint64) error {
	old, err := f.kfs.scheduler.GetFileObjectByHash(f.Hash())
	if err != nil {
		return err
	}
	content := make([]byte, size)
	copy(content, old.Content)
	return nil
}

func (f *File) Hash() string {
	return f.hash
}
