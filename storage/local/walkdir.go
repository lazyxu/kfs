package local

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Walker[T any] struct {
	ctx      context.Context
	root     string
	visitors []Visitor[T]
}

type WalkerErr struct {
	Err      error
	FilePath string
}

func NewWalker[T any](ctx context.Context, root string, visitors ...Visitor[T]) *Walker[T] {
	return &Walker[T]{
		ctx:      ctx,
		root:     root,
		visitors: visitors,
	}
}

func (c *Walker[T]) Walk(concurrent bool) (any, error) {
	root, err := filepath.Abs(c.root)
	if err != nil {
		return nil, err
	}
	ret, err := c.walk(root, concurrent)
	return ret, err
}

func (c *Walker[T]) visitorsEnter(filename string, info os.FileInfo) bool {
	for _, visitor := range c.visitors {
		if !visitor.Enter(filename, info) {
			return false
		}
	}
	return true
}

func (c *Walker[T]) visitorsExit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []T) (T, error) {
	for _, visitor := range c.visitors {
		if visitor.HasExit() {
			return visitor.Exit(ctx, filename, info, infos, rets)
		}
	}
	var t T
	return t, nil
}

func (c *Walker[T]) walk(filePath string, concurrent bool) (ret T, err error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return
	}

	var infos []os.FileInfo
	var rets []T
	defer func() {
		ret, err = c.visitorsExit(c.ctx, filePath, fileInfo, infos, rets)
	}()
	if !c.visitorsEnter(filePath, fileInfo) {
		return
	}

	if !fileInfo.IsDir() {
		return ret, filepath.SkipDir
	}
	infos, err = ioutil.ReadDir(filePath)
	if err != nil {
		return
	}

	ch := make(chan T)
	for _, info := range infos {
		select {
		case <-c.ctx.Done():
			return ret, errors.New("context deadline exceed")
		default:
			if !concurrent {
				itemFilePath := filepath.Join(filePath, info.Name())
				itemRet, _ := c.walk(itemFilePath, concurrent)
				rets = append(rets, itemRet)
				continue
			}
			go func(info os.FileInfo, concurrent bool) {
				itemFilePath := filepath.Join(filePath, info.Name())
				itemRet, err := c.walk(itemFilePath, concurrent)
				if err != nil {
					println(err.Error())
				}
				ch <- itemRet
			}(info, concurrent)
		}
	}
	if !concurrent {
		return
	}
	for i := 0; i < len(infos); i++ {
		rets = append(rets, <-ch)
	}
	return
}
