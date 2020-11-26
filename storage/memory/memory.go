package memory

import (
	"bytes"
	"crypto/sha256"
	"io"
	"sync"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/core/e"
)

type Storage struct {
	mutex sync.RWMutex
	objs  map[int]map[string][]byte
}

func New() *Storage {
	objs := make(map[int]map[string][]byte, 16)
	objs[storage.TypTree] = make(map[string][]byte, 16)
	objs[storage.TypBlob] = make(map[string][]byte, 16)
	return &Storage{
		objs: objs,
	}
}

func (s *Storage) Read(typ int, key string) (io.Reader, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	typedObjs, ok := s.objs[typ]
	if !ok {
		return nil, e.EInvalidType
	}
	data, ok := typedObjs[key]
	if !ok {
		return nil, e.ErrNotExist
	}
	return bytes.NewReader(data), nil
}

func (s *Storage) Write(typ int, reader io.Reader) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	data := buf.Bytes()

	hash := sha256.New()
	_, err = hash.Write(data)
	if err != nil {
		return "", err
	}
	key := string(hash.Sum(nil))

	typedObjs, ok := s.objs[typ]
	if !ok {
		return "", e.EInvalidType
	}
	if _, ok := typedObjs[key]; ok {
		return key, nil
	}
	typedObjs[key] = data
	return key, nil
}

func (s *Storage) Delete(typ int, key string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	typedObjs, ok := s.objs[typ]
	if !ok {
		return e.EInvalidType
	}
	delete(typedObjs, key)
	return nil
}
