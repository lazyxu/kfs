package storage

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

func (s *Storage) WriteObject(hash []byte, fn func(func(reader io.Reader))) {
	p := path.Join(s.root, "object", hex.EncodeToString(hash))
	_, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
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
}

func (s *Storage) ReadObject(hash string, fn func(reader io.Reader)) {
	p := path.Join(s.root, "object", hash)
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fn(file)
}
