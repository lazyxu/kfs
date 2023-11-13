package local_file

import (
	"sync"
)

type DriverLocalFile struct {
	driverId uint64
	srcPath  string
	encoder  string
	mutex    sync.Locker
	taskInfo TaskInfo
}
