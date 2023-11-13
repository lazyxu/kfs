package local_file

import (
	"github.com/lazyxu/kfs/core"
	"sync"
)

type DriverLocalFile struct {
	kfsCore    *core.KFS
	driverId   uint64
	DstPath    string
	Encoder    string
	Concurrent int
	mutex      sync.Locker
	taskInfo   TaskInfo
}
