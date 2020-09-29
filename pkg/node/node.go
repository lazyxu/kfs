package node

import (
	"os"
	"time"
)

type Node interface {
	Name() string
	Size() (int64, error)
	BirthTime() time.Time
	ModTime() time.Time
	ChangeTime() time.Time
	IsDir() bool
	IsFile() bool
	Mode() (mode os.FileMode)

	Truncate(size uint64) error
}
