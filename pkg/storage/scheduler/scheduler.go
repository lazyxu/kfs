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
	if hash == EmptyFileHash {
		return EmptyFile, nil
	}
	obj, err := s.storages[0].Get(hash)
	if err != nil {
		return nil, err
	}
	if file, ok := obj.(*object.File); ok {
		return file, nil
	}
	return nil, e.ENotFile
}

func (s *Scheduler) GetObjectByHash(hash string) (object.Object, error) {
	if hash == EmptyFileHash {
		return EmptyFile, nil
	}
	return s.storages[0].Get(hash)
}

func (s *Scheduler) SetObjectByHash(object object.Object) error {
	if object == EmptyFile {
		return nil
	}
	return s.storages[0].Add(object)
}

func (s *Scheduler) checkpoint() error {
	return nil
}
