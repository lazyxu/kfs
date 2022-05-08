package local

import (
	"bytes"
	"os"
	"testing"
)

const testDir = "tmp"

func TestErrorHash(t *testing.T) {
	os.RemoveAll(testDir)
	s, err := New(testDir)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.Write("error-hash", bytes.NewReader([]byte("abc")))
	if err == nil {
		t.Error("should have invalid hash")
		return
	}
	os.RemoveAll(testDir)
}

func TestWriteTwice(t *testing.T) {
	os.RemoveAll(testDir)
	s, err := New("tmp")
	if err != nil {
		t.Error(err)
		return
	}
	hash, content := NewContent("abc")
	exist, err := s.Write(hash, bytes.NewReader(content))
	if err != nil {
		t.Error(err)
		return
	}
	if exist == true {
		t.Error("should not exist")
		return
	}
	exist, err = s.Write(hash, bytes.NewReader(content))
	if err != nil {
		t.Error(err)
		return
	}
	if exist == false {
		t.Error("should exist")
		return
	}
	os.RemoveAll(testDir)
}
