package object

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/lazyxu/kfs/core/e"
	"github.com/lazyxu/kfs/scheduler"
	"github.com/lazyxu/kfs/storage"
)

type Tree struct {
	Items []*Metadata
}

var EmptyDir = &Tree{
	Items: make([]*Metadata, 0),
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

func (o *Tree) GetNode(name string) (*Metadata, error) {
	for _, it := range o.Items {
		if it.Name == name {
			return it, nil
		}
	}
	return nil, e.ErrNotExist
}

func (o *Tree) Write(s *scheduler.Scheduler) (string, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(o)
	if err != nil {
		return "", e.EWriteObject
	}
	return s.WriteStream(storage.TypDir, &b)
}

func (o *Tree) Read(s *scheduler.Scheduler, key string) error {
	reader, err := s.ReadStream(storage.TypDir, key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(reader).Decode(o)
}

func ReadDir(s *scheduler.Scheduler, key string) (*Tree, error) {
	d := new(Tree)
	err := d.Read(s, key)
	return d, err
}

func (o *Tree) IsDir() bool {
	return true
}

func (o *Tree) IsFile() bool {
	return false
}
