package local

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	StateNew = iota
	StateScan
	StateUpload
	StateDone
	StateStop
)

type BackupCtx[T any] struct {
	ctx            context.Context
	root           string
	mutex          sync.RWMutex
	done           bool
	canceled       bool
	errs           []BackupErr
	scanProcess    []int
	scanMaxProcess []int
	curFilename    []int
	visitors       []Visitor[T]
}

type BackupErr struct {
	Err      error
	FilePath string
}

func NewBackupCtx[T any](ctx context.Context, root string, visitors ...Visitor[T]) *BackupCtx[T] {
	return &BackupCtx[T]{
		ctx:      ctx,
		root:     root,
		errs:     []BackupErr{},
		visitors: visitors,
	}
}

func (c *BackupCtx[T]) Scan() (any, error) {
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

func (c *BackupCtx[T]) visitorsEnter(filename string, info os.FileInfo) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, visitor := range c.visitors {
		if !visitor.Enter(filename, info) {
			return false
		}
	}
	return true
}

func (c *BackupCtx[T]) visitorsExit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []T) (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, visitor := range c.visitors {
		if visitor.HasExit() {
			return visitor.Exit(ctx, filename, info, infos, rets)
		}
	}
	var t T
	return t, nil
}

func (c *BackupCtx[T]) walk(filename string) (ret T, err error) {
	info, err := os.Lstat(filename)
	if err != nil {
		return
	}

	var infos []os.FileInfo
	var rets []T
	defer func() {
		ret, err = c.visitorsExit(c.ctx, filename, info, infos, rets)
	}()
	if !c.visitorsEnter(filename, info) {
		return
	}

	if !info.IsDir() {
		return ret, filepath.SkipDir
	}
	infos, err = ioutil.ReadDir(filename)
	if err != nil {
		return
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
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.scanProcess = nil
			c.mutex.Unlock()
			return ret, errors.New("context deadline exceed")
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
	return
}
