package local_file

import (
	"sync"
)

type DriverLocalFile struct {
	driverId uint64
	mutex    sync.Locker
	taskInfo TaskInfo
}
