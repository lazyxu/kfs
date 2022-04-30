package storage

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type ObjectType int

const (
	ObjectTypeFile ObjectType = iota
	ObjectTypeDirectory
	ObjectTyeCommit
)

type Object interface {
	Type() ObjectType
}

func (s *Storage) WriteObject(hash []byte, fn func(func(reader io.Reader))) {
	s.writeObject(hex.EncodeToString(hash), func(f *os.File) {
		hw := s.HashFunc()
		fn(func(reader io.Reader) {
			rr := io.TeeReader(reader, hw)
			_, err := io.Copy(f, rr)
			if err != nil {
				panic(err)
			}
		})
		actual := hw.Cal(nil)
		if bytes.Compare(hash, actual) != 0 {
			panic(fmt.Errorf("invalid hash: expected: %s, actual: %s", hex.EncodeToString(hash), hex.EncodeToString(actual)))
		}
	})
}

func (s *Storage) UpdateBranchHash(branch string, path string, metadata *Metadata) {
	var hash string
	s.readBranch(branch, func(b *Branch) {
		hash = b.BranchHash
	})
	dirs := make([]*Directory, 0)
	d := s.ReadDirectory(hash)
	dirs = append(dirs, d)
loopDir:
	for _, dirname := range strings.Split(path, "/") {
		if dirname == "" {
			panic(errors.New("非法路径名"))
		}
		for _, item := range *d {
			if item.Name == dirname {
				hash = item.Hash
				continue loopDir
			}
		}
		d := s.ReadDirectory(hash)
		dirs = append(dirs, d)
	}
	for i := len(dirs) - 1; i > 0; i-- {
		cur := dirs[i]
		cur = append(*cur, metadata)
		s.WriteDirectory(d)
	}
}

func (s *Storage) GetObjectReader(hash string, fn func(reader io.Reader)) {
	p := path.Join(s.root, "object", hash)
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fn(file)
}

func (s *Storage) ReadDirectory(hash string) *Directory {
	p := path.Join(s.root, "object", hash)
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	i := &Directory{}
	DefaultDirectoryEncoderDecoder.Decode(i, file)
	return i
}

func (s *Storage) WriteDirectory(d *Directory) {
	s.writeObject(hash, func(f *os.File) {
		hw := s.HashFunc()
		fn(func(reader io.Reader) {
			rr := io.TeeReader(reader, hw)
			_, err := io.Copy(f, rr)
			if err != nil {
				panic(err)
			}
		})
		actual := hw.Cal(nil)
		if bytes.Compare(hash, actual) != 0 {
			panic(fmt.Errorf("invalid hash: expected: %s, actual: %s", hex.EncodeToString(hash), hex.EncodeToString(actual)))
		}
	})
}
