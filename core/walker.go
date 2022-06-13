package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type File[T any] struct {
	Path     string
	Info     os.FileInfo
	parent   chan T
	children chan T
}

type WorkHandlers[T any] interface {
	FilePathFilter(filePath string) bool
	FileInfoFilter(filePath string, info os.FileInfo) bool
	FileHandler(ctx context.Context, filePath string, info os.FileInfo, children []T) T
	ErrHandler(filePath string, err error)
	StackSizeHandler(size int)
}

type DefaultWalkHandlers[T any] struct{}

func (DefaultWalkHandlers[T]) FilePathFilter(filePath string) bool {
	return false
}

func (DefaultWalkHandlers[T]) FileInfoFilter(filePath string, info os.FileInfo) bool {
	return false
}

func (DefaultWalkHandlers[T]) FileHandler(ctx context.Context, filePath string, info os.FileInfo, children []T) (t T) {
	return
}

func (DefaultWalkHandlers[T]) ErrHandler(filePath string, err error) {
	println(filePath, err.Error())
}

func (DefaultWalkHandlers[T]) StackSizeHandler(size int) {
}

func Walk[T any](ctx context.Context, filePath string, concurrent int, handlers WorkHandlers[T]) (t T, err error) {
	if concurrent <= 0 {
		return t, fmt.Errorf("concurrent should be > 0, actual %d", concurrent)
	}
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return
	}
	stack := arraystack.New()
	stack.Push(&File[T]{
		Path: filePath,
	})
	ch := make(chan struct{}, concurrent)
	defer func() {
		for i := 0; i < len(ch); i++ {
			ch <- struct{}{}
		}
	}()
	for !stack.Empty() {
		select {
		case <-ctx.Done():
			return t, context.DeadlineExceeded
		default:
		}
		v, _ := stack.Pop()
		if v != nil {
			handlers.StackSizeHandler(stack.Size())
			f := v.(*File[T])
			var info os.FileInfo
			info, continues := preHandler(f.Path, handlers)
			if continues {
				if f.parent != nil {
					f.parent <- t
				}
				continue
			}
			cur := &File[T]{
				Path:   f.Path,
				Info:   info,
				parent: f.parent,
			}
			stack.Push(cur)
			stack.Push(nil)

			if !info.IsDir() {
				continue
			}

			var infos []os.DirEntry
			infos, err = os.ReadDir(f.Path)
			if err != nil {
				handlers.ErrHandler(f.Path, err)
				continue
			}
			cur.children = make(chan T, len(infos))
			for i := len(infos) - 1; i >= 0; i-- {
				stack.Push(&File[T]{
					Path:   filepath.Join(f.Path, infos[i].Name()),
					parent: cur.children,
				})
			}
		} else {
			vv, ok := stack.Pop()
			if !ok {
				panic(errors.New("non-nil element was pushed into stack before nil"))
			}
			handlers.StackSizeHandler(stack.Size())
			f := vv.(*File[T])
			ch <- struct{}{}
			go func() {
				cnt := cap(f.children)
				children := make([]T, cnt)
				for i := 0; i < cnt; i++ {
					children[i] = <-f.children
				}
				parent := handlers.FileHandler(ctx, f.Path, f.Info, children)
				if f.parent != nil {
					f.parent <- parent
				} else {
					t = parent
				}
				<-ch
			}()
		}
	}
	handlers.StackSizeHandler(0)
	return
}

func preHandler[T any](filePath string, handlers WorkHandlers[T]) (info os.FileInfo, continues bool) {
	var err error
	defer func() {
		if err != nil {
			handlers.ErrHandler(filePath, err)
			continues = true
		}
	}()
	if handlers.FilePathFilter(filePath) {
		return info, true
	}
	info, err = os.Lstat(filePath)
	if err != nil {
		return
	}
	if handlers.FileInfoFilter(filePath, info) {
		return info, true
	}
	return
}
