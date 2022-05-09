package local

import (
	"bytes"
	"os"
	"path"
	"sync"
	"testing"
)

const testDir = "tmp"

func TestNew(t *testing.T) {
	os.RemoveAll(testDir)
	_, err := New(testDir)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = New(testDir)
	if err != nil {
		t.Error(err)
		return
	}
	err = os.Remove(path.Join(testDir, lockFileName))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = New(testDir)
	if err != nil {
		t.Error(err)
		return
	}
}

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
}

func TestConcurrentWrite(t *testing.T) {
	os.RemoveAll(testDir)
	s, err := New("tmp")
	if err != nil {
		t.Error(err)
		return
	}
	hash, content := NewContent("abc")
	errCnt := 0
	existCnt := 0
	var wg sync.WaitGroup
	n := 10
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			exist, err := s.Write(hash, bytes.NewReader(content))
			if err != nil {
				errCnt += 1
			}
			if exist {
				existCnt += 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if errCnt != 0 {
		t.Error("errCnt should be 0")
		return
	}
	if existCnt != n-1 {
		t.Errorf("existCnt should be %d", n-1)
		return
	}
}
