package local

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	StateNew = iota
	StateScan
	StateUpload
	StateDone
	StateStop
)

type BackupCtx struct {
	ctx            context.Context
	root           string
	mutex          sync.RWMutex
	done           bool
	canceled       bool
	errs           []BackupErr
	scanProcess    []int
	scanMaxProcess []int
	curFilename    []int
	visitors       []Visitor
}

type BackupErr struct {
	Err      error
	FilePath string
}

func NewBackupCtx(ctx context.Context, root string, visitors ...Visitor) *BackupCtx {
	return &BackupCtx{
		ctx:      ctx,
		root:     root,
		errs:     []BackupErr{},
		visitors: visitors,
	}
}

func (c *BackupCtx) Scan() (any, error) {
	defer func() {
		c.mutex.Lock()
		c.done = true
		c.mutex.Unlock()
	}()
	root, err := filepath.Abs(c.root)
	if err != nil {
		return nil, err
	}
	ret, err := c.walk(root)
	return ret, err
}

func (c *BackupCtx) visitorsEnter(filename string, info os.FileInfo) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, visitor := range c.visitors {
		if !visitor.Enter(filename, info) {
			return false
		}
	}
	return true
}

func (c *BackupCtx) visitorsExit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []any) (any, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, visitor := range c.visitors {
		if visitor.HasExit() {
			return visitor.Exit(ctx, filename, info, infos, rets)
		}
	}
	return nil, nil
}

func (c *BackupCtx) walk(filename string) (ret any, err error) {

	info, err := os.Lstat(filename)
	if err != nil {
		return nil, err
	}

	var infos []os.FileInfo
	var rets []any
	defer func() {
		ret, err = c.visitorsExit(c.ctx, filename, info, infos, rets)
	}()
	if !c.visitorsEnter(filename, info) {
		return nil, nil
	}

	if !info.IsDir() {
		return nil, filepath.SkipDir
	}
	infos, err = ioutil.ReadDir(filename)
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	if len(infos) != 0 {
		c.scanProcess = append(c.scanProcess, 0)
		c.scanMaxProcess = append(c.scanMaxProcess, len(infos))
	}
	c.mutex.Unlock()
	defer func() {
		c.mutex.Lock()
		if len(infos) != 0 {
			c.scanProcess = c.scanProcess[0 : len(c.scanProcess)-1]
			c.scanMaxProcess = c.scanMaxProcess[0 : len(c.scanMaxProcess)-1]
		}
		c.mutex.Unlock()
	}()

	for _, info := range infos {
		time.Sleep(time.Millisecond * 9)
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.scanProcess = nil
			c.mutex.Unlock()
			return nil, errors.New("context deadline exceed")
		default:
			filename := filepath.Join(filename, info.Name())
			ret, err := c.walk(filename)
			c.mutex.Lock()
			rets = append(rets, ret)
			c.mutex.Unlock()
			if err == filepath.SkipDir {
				c.mutex.Lock()
				c.scanProcess[len(c.scanProcess)-1]++
				c.mutex.Unlock()
				continue
			}
			if err != nil {
				c.mutex.Lock()
				c.errs = append(c.errs, BackupErr{
					Err:      err,
					FilePath: filename,
				})
				c.scanProcess[len(c.scanProcess)-1]++
				c.mutex.Unlock()
				continue
			}
		}
		c.mutex.Lock()
		c.scanProcess[len(c.scanProcess)-1]++
		c.mutex.Unlock()
	}
	return nil, nil
}
