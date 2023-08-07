package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []T) T
	OnFileError(filePath string, info os.FileInfo, err error)
	StackSizeHandler(size int)
	StartWorker(ctx context.Context, index int)
	EndWorker(ctx context.Context, index int)
	PushFile(info os.FileInfo)
	StartFile(ctx context.Context, index int, filePath string, info os.FileInfo)
	HasPushedAllToStack()
}

type DefaultWalkHandlers[T any] struct{}

func (DefaultWalkHandlers[T]) FilePathFilter(filePath string) bool {
	return false
}

func (DefaultWalkHandlers[T]) FileInfoFilter(filePath string, info os.FileInfo) bool {
	return false
}

func (DefaultWalkHandlers[T]) StartFile(ctx context.Context, index int, filePath string, info os.FileInfo) {
}

func (DefaultWalkHandlers[T]) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []T) (t T) {
	return
}

func (DefaultWalkHandlers[T]) OnFileError(filePath string, info os.FileInfo, err error) {
	println(filePath, err.Error())
}

func (DefaultWalkHandlers[T]) StackSizeHandler(size int) {
}

func (DefaultWalkHandlers[T]) StartWorker(ctx context.Context, index int) {
}

func (DefaultWalkHandlers[T]) EndWorker(ctx context.Context, index int) {
}

func (DefaultWalkHandlers[T]) PushFile(info os.FileInfo) {
}
func (DefaultWalkHandlers[T]) HasPushedAllToStack() {
}

func Walk[T any](ctx context.Context, filePath string, concurrent int, handlers WorkHandlers[T]) (t T, err1 error) {
	if concurrent <= 0 {
		return t, fmt.Errorf("concurrent should be > 0, actual %d", concurrent)
	}
	filePath, err1 = filepath.Abs(filePath)
	if err1 != nil {
		return
	}
	stack := arraystack.New()
	stack.Push(&File[T]{
		Path: filePath,
	})
	ch := make(chan *File[T], 1000000)
	var wg sync.WaitGroup
	defer func() {
		close(ch)
		wg.Wait()
	}()
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func(index int) {
			handlers.StartWorker(ctx, index)
			for {
				select {
				case <-ctx.Done():
					err1 = context.DeadlineExceeded
					break
				default:
				}
				f, ok := <-ch
				if !ok {
					break
				}
				handlers.StartFile(ctx, index, f.Path, f.Info)
				cnt := cap(f.children)
				children := make([]T, cnt)
				for i := 0; i < cnt; i++ {
					select {
					case children[i] = <-f.children:
					case <-ctx.Done():
						err1 = context.DeadlineExceeded
						goto EndWorker
					}
				}
				parent := handlers.FileHandler(ctx, index, f.Path, f.Info, children)
				if f.parent != nil {
					select {
					case f.parent <- parent:
					case <-ctx.Done():
						err1 = context.DeadlineExceeded
						goto EndWorker
					}

				} else {
					t = parent
				}
			}
		EndWorker:
			handlers.EndWorker(ctx, index)
			wg.Done()
		}(i)
	}
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

			infos, err := os.ReadDir(f.Path)
			if err != nil {
				handlers.OnFileError(f.Path, f.Info, err)
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
			handlers.PushFile(f.Info)
			ch <- f
		}
	}
	handlers.StackSizeHandler(0)
	handlers.HasPushedAllToStack()
	return
}

func preHandler[T any](filePath string, handlers WorkHandlers[T]) (info os.FileInfo, continues bool) {
	var err error
	defer func() {
		if err != nil {
			handlers.OnFileError(filePath, info, err)
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
