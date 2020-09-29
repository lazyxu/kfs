package memory

import (
	"sync"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/kfs/e"
)

type Storage struct {
	mutex   sync.RWMutex
	objects map[string]object.Object
}

func New() *Storage {
	return &Storage{
		objects: make(map[string]object.Object, 16),
	}
}

func (s *Storage) Get(hash string) (object.Object, error) {
	s.mutex.RLock()
	n, ok := s.objects[hash]
	s.mutex.RUnlock()
	if !ok {
		return nil, e.ErrNotExist
	}
	return n, nil
}

func (s *Storage) Add(object object.Object) error {
	s.mutex.Lock()
	s.objects[object.Hash()] = object
	s.mutex.Unlock()
	return nil
}

func (s *Storage) Remove(hash string) error {
	s.mutex.Lock()
	delete(s.objects, hash)
	s.mutex.Unlock()
	return nil
}
