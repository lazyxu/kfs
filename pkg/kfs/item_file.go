package kfs

import (
	"fmt"
	"os"

	"github.com/lazyxu/kfs/object"
	"github.com/sirupsen/logrus"
)

type ItemFile struct {
	DItem
	object *object.File
}

func newFile(kfs *KFS, name string) *ItemFile {
	return &ItemFile{
		DItem: DItem{
			kfs: kfs,
		},
		object: object.NewEmptyFile(name),
	}
}

func (i *ItemFile) Name() string {
	return i.object.Name()
}

func (i *ItemFile) IsDir() bool {
	return i.object.IsDir()
}

func (i *ItemFile) IsFile() bool {
	return i.object.IsFile()
}

func (i *ItemFile) loadObject() error {
	if i.object != nil {
		return nil
	}
	obj, err := i.kfs.scheduler.GetFileObjectByHash(i.object.Hash())
	if err != nil {
		return err
	}
	i.object = obj
	return nil
}

func (i *ItemFile) GetContent() (string, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	err := i.loadObject()
	if err != nil {
		return "", err
	}
	return i.object.Content(), nil
}

func (i *ItemFile) Size() (int64, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	err := i.loadObject()
	if err != nil {
		return 0, err
	}
	return int64(len(i.object.Content())), nil
}

func (i *ItemFile) SetContent(content []byte, offset int64) (int64, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"offset":  offset,
		"len":     len(content),
	}).Debug("SetContent")
	err := i.loadObject()
	if err != nil {
		return 0, err
	}
	fmt.Println("1")
	content = append([]byte(i.object.Content())[0:offset], content...)
	fmt.Println("2")
	logrus.WithFields(logrus.Fields{
		"content": string(content),
		"len":     len(content),
		"name":    i.Name(),
	}).Debug("AfterSetContent")
	fmt.Println("3")
	i.object = i.object.Clone().(*object.File)
	i.object.SetContent(string(content))
	i.object.SetHash(i.object.Hash())
	return int64(len(i.object.Content())), nil
}

func (i *ItemFile) Mode() os.FileMode {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.object.Mode()
}

func (i *ItemFile) Truncate(size uint64) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err := i.loadObject()
	if err != nil {
		return err
	}
	content := make([]byte, size)
	copy(content, i.object.Content())
	i.object = i.object.Clone().(*object.File)
	i.object.SetContent(string(content))
	i.object.SetHash(i.object.Hash())
	return nil
}

func (i *ItemFile) Item() object.Object {
	return i.object
}
