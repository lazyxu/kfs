package memory

import (
	"bytes"
	"io"
	"sync"

	"github.com/lazyxu/kfs/kfscore/kfscrypto"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/lazyxu/kfs/kfscore/e"
)

type Storage struct {
	storage.BaseStorage
	mutex    sync.RWMutex
	objs     map[int]map[string][]byte
	tempObjs map[int]map[string][]byte
	refs     map[string]string
}

func New(hashFunc func() kfscrypto.Hash) *Storage {
	objs := make(map[int]map[string][]byte, 16)
	objs[storage.TypTree] = make(map[string][]byte, 16)
	objs[storage.TypBlob] = make(map[string][]byte, 16)
	tempObjs := make(map[int]map[string][]byte, 16)
	tempObjs[storage.TypTree] = make(map[string][]byte, 16)
	tempObjs[storage.TypBlob] = make(map[string][]byte, 16)
	return &Storage{
		BaseStorage: storage.NewBase(hashFunc),
		objs:        objs,
		tempObjs:    tempObjs,
		refs:        make(map[string]string, 8),
	}
}

func (s *Storage) Read(typ int, key string, f func(reader io.Reader) error) error {
	s.mutex.RLock()
	typedObjs, ok := s.objs[typ]
	if !ok {
		s.mutex.RUnlock()
		return e.EInvalidType
	}
	data, ok := typedObjs[key]
	s.mutex.RUnlock()
	if !ok {
		return e.ErrNotExist
	}
	return f(bytes.NewReader(data))
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

func (s *Storage) Commit(typ int, key string) error {
	return nil
}

func (s *Storage) Exist(typ int, key string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, exist := s.objs[typ]
	return exist, nil
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

func (s *Storage) GetRefs() ([]string, error) {
	s.mutex.RLock()
	var branches []string
	for name := range s.refs {
		branches = append(branches, name)
	}
	s.mutex.RUnlock()
	return branches, nil
}
func (s *Storage) Status() (status storage.Status, err error) {
	s.mutex.RLock()
	for _, i := range s.objs {
		for _, j := range i {
			status.TotalPhysicalSize += uint64(len(j))
		}
	}
	for _, i := range s.objs[storage.TypBlob] {
		status.BlobLogicalSize += uint64(len(i))
	}
	status.BlobCount = uint64(len(s.objs[storage.TypBlob]))
	status.TreeCount = uint64(len(s.objs[storage.TypTree]))
	s.mutex.RUnlock()
	return status, err
}
