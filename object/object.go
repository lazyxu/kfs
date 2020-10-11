package object

import "github.com/lazyxu/kfs/scheduler"

type Object interface {
	IsDir() bool
	IsFile() bool
	Write(s *scheduler.Scheduler) (string, error)
	Read(s *scheduler.Scheduler, key string) error
}
