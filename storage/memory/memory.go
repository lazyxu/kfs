package memory

import (
	"bytes"
	"io"
	"sync"

	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/core/e"
)

type Storage struct {
	storage.BaseStorage
	mutex sync.RWMutex
	objs  map[int]map[string][]byte
	refs  map[string]string
}

func New(hashFunc func() kfscrypto.Hash, checkOnWrite bool, checkOnRead bool) *Storage {
	objs := make(map[int]map[string][]byte, 16)
	objs[storage.TypTree] = make(map[string][]byte, 16)
	objs[storage.TypBlob] = make(map[string][]byte, 16)
	return &Storage{
		BaseStorage: storage.NewBase(hashFunc, checkOnWrite, checkOnRead),
		objs:        objs,
		refs:        make(map[string]string, 8),
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
	hw := s.HashFunc()
	rr := io.TeeReader(reader, hw)
	_, err := buf.ReadFrom(rr)
	if err != nil {
		return "", err
	}
	data := buf.Bytes()

	key, err := hw.Cal(nil)
	if err != nil {
		return "", err
	}

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

func (s *Storage) UpdateRef(name string, expect string, desire string) error {
	s.mutex.Lock()
	if expect != "" && s.refs[name] != expect {
		return e.ErrInvalid
	}
	s.refs[name] = desire
	s.mutex.Unlock()
	return nil
}

func (s *Storage) GetRef(name string) (string, error) {
	s.mutex.RLock()
	hash, ok := s.refs[name]
	s.mutex.RUnlock()
	if !ok {
		return "", e.ErrNotExist
	}
	return hash, nil
}
