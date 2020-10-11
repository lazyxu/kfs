package obj

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/lazyxu/kfs/kfs/e"
	"github.com/lazyxu/kfs/storage"
	"github.com/lazyxu/kfs/storage/scheduler"
)

type Dir struct {
	Items []Metadata
}

var EmptyDir = &Dir{
	Items: make([]Metadata, 0),
}
var EmptyDirHash string

func init() {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(EmptyDir)
	if err != nil {
		panic(err)
	}
	hash := sha256.New()
	_, err = hash.Write(b.Bytes())
	EmptyDirHash = string(hash.Sum(nil))
}

func (o *Dir) GetNode(name string) (*Metadata, error) {
	for _, it := range o.Items {
		if it.Name == name {
			return &it, nil
		}
	}
	return nil, e.ErrNotExist
}

func (o *Dir) Write(s *scheduler.Scheduler) (string, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(o)
	if err != nil {
		return "", e.EWriteObject
	}
	return s.WriteStream(storage.TypDir, &b)
}

func (o *Dir) Read(s *scheduler.Scheduler, key string) error {
	reader, err := s.ReadStream(storage.TypDir, key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(reader).Decode(o)
}

func ReadDir(s *scheduler.Scheduler, key string) (*Dir, error) {
	d := new(Dir)
	err := d.Read(s, key)
	return d, err
}

func (o *Dir) IsDir() bool {
	return true
}

func (o *Dir) IsFile() bool {
	return false
}
