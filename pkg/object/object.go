package object

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

type Object interface {
	Hash() string
}

type File struct {
	Content string
}

func (f *File) Hash() string {
	return string(sha256.New().Sum([]byte(f.Content)))
}

type Dir struct {
	Items map[string]*Item
}

type Item struct {
	obj   Object
	name  string
	BTime time.Time
	Mtime time.Time
}

func (d *Dir) Hash() string {
	hash := sha256.New()
	for _, item := range d.Items {
		binary.Write(hash, binary.LittleEndian, item.name)
		binary.Write(hash, binary.LittleEndian, item.BTime.UnixNano())
		binary.Write(hash, binary.LittleEndian, item.Mtime.UnixNano())
	}
	return string(hash.Sum(nil))
}
