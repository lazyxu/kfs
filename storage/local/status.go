package local

import (
	"strconv"
	"strings"
)

type Status struct {
	Done        bool
	Canceled    bool
	Errs        []BackupErr
	ScanProcess string
}

func (c *BackupCtx) GetStatus() Status {
	c.mutex.RLock()
	var scanProcess string
	if len(c.scanProcess) == 0 {
		scanProcess = "已完成"
	} else {
		process := new(strings.Builder)
		for i, p := range c.scanProcess {
			if i != 0 {
				process.WriteByte(',')
			}
			process.WriteString(strconv.Itoa(p))
			process.WriteString("/")
			process.WriteString(strconv.Itoa(c.scanMaxProcess[i]))
		}
		scanProcess = process.String()
	}
	p := Status{
		Canceled:    c.canceled,
		Errs:        c.errs,
		ScanProcess: scanProcess,
		Done:        c.done,
	}
	c.mutex.RUnlock()
	return p
}
