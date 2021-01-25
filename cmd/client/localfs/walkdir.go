package localfs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/dustin/go-humanize"
)

type BackUpCtx struct {
	ctx          context.Context
	root         string
	mutex        sync.RWMutex
	fileSize     uint64
	fileCount    uint64
	dirCount     uint64
	largeFiles   map[string]interface{}
	ignoredFiles []string
	ignoreRules  []IgnoreRule
	done         bool
	canceled     bool
	errs         []BackUpErr
}

type BackUpErr struct {
	Err      error
	FileName string
}

func NewBackUpCtx(ctx context.Context, root string, ignoreRules []IgnoreRule) *BackUpCtx {
	return &BackUpCtx{
		ctx:          ctx,
		root:         root,
		largeFiles:   make(map[string]interface{}),
		ignoredFiles: []string{},
		ignoreRules:  ignoreRules,
		errs:         []BackUpErr{},
	}
}

func (c *BackUpCtx) Scan() {
	c.walk(c.root)
	c.mutex.Lock()
	c.done = true
	c.mutex.Unlock()
}

type Status struct {
	FileSize     string
	FileCount    uint64
	DirCount     uint64
	LargeFiles   map[string]interface{}
	IgnoredFiles []string
	Done         bool
	Canceled     bool
	Errs         []BackUpErr
}

func (c *BackUpCtx) GetStatus() Status {
	c.mutex.RLock()
	p := Status{
		FileSize:     humanize.Bytes(c.fileSize),
		FileCount:    c.fileCount,
		DirCount:     c.dirCount,
		LargeFiles:   c.largeFiles,
		IgnoredFiles: c.ignoredFiles,
		Done:         c.done,
		Canceled:     c.canceled,
		Errs:         c.errs,
	}
	c.mutex.RUnlock()
	return p
}

func (c *BackUpCtx) ignoreFile(fileName string) bool {
	if ignoreByStd(fileName) {
		return true
	}
	for _, rule := range c.ignoreRules {
		if rule.Ignore(fileName) {
			return true
		}
	}
	return false
}

func (c *BackUpCtx) walk(fileName string) {
	info, err := os.Lstat(fileName)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FileName: fileName,
		})
		c.mutex.Unlock()
		return
	}
	modeType := info.Mode() & os.ModeType
	if c.ignoreFile(fileName) {
		c.mutex.Lock()
		c.ignoredFiles = append(c.ignoredFiles, fileName)
		c.mutex.Unlock()
		return
	}
	if modeType == 0 {
		c.mutex.Lock()
		c.fileCount++
		c.fileSize += uint64(info.Size())
		if info.Size() > 100*1024*1024 {
			c.largeFiles[fileName] = humanize.Bytes(uint64(info.Size()))
		}
		c.mutex.Unlock()
		return
	}
	if modeType&os.ModeSymlink != 0 {
		c.mutex.Lock()
		c.fileCount++
		c.mutex.Unlock()
		return
	}
	if !info.IsDir() {
		return
	}
	infos, err := ioutil.ReadDir(fileName)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FileName: fileName,
		})
		c.mutex.Unlock()
		return
	}
	c.mutex.Lock()
	c.dirCount += 1
	c.mutex.Unlock()

	for _, info := range infos {
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.mutex.Unlock()
			return
		default:
			filename := filepath.Join(fileName, info.Name())
			c.walk(filename)
		}
	}
	return
}
