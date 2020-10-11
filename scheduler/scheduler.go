package scheduler

import (
	"io"

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

func (s *Scheduler) checkpoint() error {
	return nil
}

func (s *Scheduler) WriteStream(typ int, reader io.Reader) (string, error) {
	return s.storages[0].Write(typ, reader)
}

func (s *Scheduler) ReadStream(typ int, key string) (io.Reader, error) {
	return s.storages[0].Read(typ, key)
}

func (s *Scheduler) DeleteStream(typ int, key string) error {
	return s.storages[0].Delete(typ, key)
}
