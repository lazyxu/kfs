package storage

import (
	"os"
	"path"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

func (s *Storage) readFile(p string, cb func(f *os.File)) {
	f, err := os.OpenFile(p, os.O_RDONLY, filePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cb(f)
}

func (s *Storage) writeObject(hash string, cb func(f *os.File)) {
	p := path.Join(s.root, "object", hash)
	_, err := os.Stat(p)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		panic(err)
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cb(f)
}
