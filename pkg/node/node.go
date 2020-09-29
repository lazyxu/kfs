package node

import (
	"os"
	"time"
)

type Node interface {
	Name() string
	Size() (int64, error)
	BirthTime() time.Time
	AccessTime() time.Time
	ModTime() time.Time
	ChangeTime() time.Time
	IsDir() bool
	IsFile() bool
	Mode() (mode os.FileMode)

	Truncate(size uint64) error
}

type TimeImpl struct {
	BTime time.Time
	ATime time.Time
	Mtime time.Time
	CTime time.Time
}

func (t *TimeImpl) BirthTime() time.Time {
	return t.BTime
}
func (t *TimeImpl) AccessTime() time.Time {
	return t.ATime
}
func (t *TimeImpl) ModTime() time.Time {
	return t.Mtime
}
func (t *TimeImpl) ChangeTime() time.Time {
	return t.CTime
}
