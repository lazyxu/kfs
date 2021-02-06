package localfs

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

type Status struct {
	FileSize       string
	FileCount      uint64
	DirCount       uint64
	LargeFiles     map[string]interface{}
	IgnoredFiles   []string
	Done           bool
	Canceled       bool
	Errs           []BackUpErr
	ScanProcess    string
	UploadingCount int
}

func (c *BackUpCtx) GetStatus() Status {
	c.mutex.RLock()
	var scanProcess string
	if len(c.scanProcess) == 0 {
		scanProcess = "已完成"
	} else {
		process := new(strings.Builder)
		for i, p := range c.scanProcess {
			if i != 0 {
				process.WriteByte('.')
			}
			process.WriteString(strconv.Itoa(p))
		}
		scanProcess = process.String()
	}
	p := Status{
		FileSize:       humanize.Bytes(c.fileSize),
		FileCount:      c.fileCount,
		DirCount:       c.dirCount,
		LargeFiles:     c.largeFiles,
		IgnoredFiles:   c.ignoredFiles,
		Canceled:       c.canceled,
		Errs:           c.errs,
		ScanProcess:    scanProcess,
		Done:           c.done,
		UploadingCount: c.queue.UploadingCount(),
	}
	c.mutex.RUnlock()
	return p
}
