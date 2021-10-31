package storage

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

func (s *Storage) WriteObject(hash []byte, fn func(func(reader io.Reader) error) error) error {
	p := path.Join(s.root, "object", hex.EncodeToString(hash))
	_, err := os.Stat(p)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		f.Close()
		return err
	}
	hw := s.HashFunc()
	err = fn(func(reader io.Reader) error {
		rr := io.TeeReader(reader, hw)
		_, err := io.Copy(f, rr)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	f.Close()
	actual, err := hw.Cal(nil)
	if err != nil {
		os.Remove(p)
		return err
	}
	if bytes.Equal(hash, actual) {
		return fmt.Errorf("invalid hash: expected: %s, actual: %s", hash, actual)
	}
	return nil
}

func (s *Storage) ReadObject(hash []byte, fn func(reader io.Reader) error) error {
	p := path.Join(s.root, "object", hex.EncodeToString(hash))
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return fn(file)
}
