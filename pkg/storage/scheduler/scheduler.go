package scheduler

import (
	"github.com/lazyxu/kfs/kfs/e"
	"github.com/lazyxu/kfs/object"
	"github.com/lazyxu/kfs/storage"
)

type Scheduler struct {
	storages []storage.Storage
}

func New(storages ...storage.Storage) *Scheduler {
	return &Scheduler{
		storages: storages,
	}
}

func (s *Scheduler) GetFileObjectByHash(hash string) (*object.File, error) {
	if hash == object.EmptyFileHash {
		return object.EmptyFile, nil
	}
	obj, err := s.storages[0].Get(hash)
	if err != nil {
		return nil, err
	}
	if o, ok := obj.(*object.File); ok {
		return o, nil
	}
	return nil, e.ENotFile
}

func (s *Scheduler) GetDirObjectByHash(hash string) (*object.Dir, error) {
	if hash == object.EmptyDirHash {
		return object.EmptyDir, nil
	}
	obj, err := s.storages[0].Get(hash)
	if err != nil {
		return nil, err
	}
	if o, ok := obj.(*object.Dir); ok {
		return o, nil
	}
	return nil, e.ENotFile
}

func (s *Scheduler) GetObjectByHash(hash string) (object.Object, error) {
	if hash == object.EmptyFileHash {
		return object.EmptyFile, nil
	}
	return s.storages[0].Get(hash)
}

func (s *Scheduler) SetObjectByHash(o object.Object) error {
	if o == object.EmptyFile {
		return nil
	}
	return s.storages[0].Add(o)
}

func (s *Scheduler) checkpoint() error {
	return nil
}
