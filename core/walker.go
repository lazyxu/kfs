package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type file[T any] struct {
	path     string
	info     os.FileInfo
	parent   chan T
	children chan T
}

func Walker[T any](ctx context.Context, root string, concurrent int, handle func(ctx context.Context, w *file[T]),
	errHandle func(filePath string, err error)) error {
	if concurrent <= 0 {
		return fmt.Errorf("concurrent should be > 0, actual %d", concurrent)
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	stack := arraystack.New()
	stack.Push(&file[T]{
		path: root,
	})
	ch := make(chan struct{}, concurrent)
	for !stack.Empty() {
		select {
		case <-ctx.Done():
			err = context.DeadlineExceeded
			goto exit
		default:
		}
		v, ok := stack.Pop()
		if !ok {
			panic(errors.New("stack is not empty"))
		}
		if v != nil {
			f := v.(*file[T])
			var info os.FileInfo
			info, err = os.Lstat(f.path)
			if err != nil {
				errHandle(f.path, err)
				continue
			}
			cur := &file[T]{
				path:   f.path,
				info:   info,
				parent: f.parent,
			}
			stack.Push(cur)
			stack.Push(nil)

			if !info.IsDir() {
				continue
			}

			var infos []os.DirEntry
			infos, err = os.ReadDir(f.path)
			if err != nil {
				errHandle(f.path, err)
				continue
			}
			cur.children = make(chan T, len(infos))
			for i := len(infos) - 1; i >= 0; i-- {
				stack.Push(&file[T]{
					path:   filepath.Join(f.path, infos[i].Name()),
					parent: cur.children,
				})
			}
		} else {
			vv, ok := stack.Pop()
			if !ok {
				panic(errors.New("non-nil element was pushed into stack before nil"))
			}
			f := vv.(*file[T])
			ch <- struct{}{}
			go func() {
				handle(ctx, f)
				<-ch
			}()
		}
	}
exit:
	for i := 0; i < concurrent; i++ {
		ch <- struct{}{}
	}
	return err
}
