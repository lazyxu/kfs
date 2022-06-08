package core

import (
	"context"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWalker(t *testing.T) {
	ctx := context.Background()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Error(err)
		return
	}
	err = Walker[int64](ctx, root, 15, func(ctx context.Context, f *file[int64]) {
		var size int64
		if !f.info.IsDir() {
			size = f.info.Size()
		}
		if f.children != nil {
			for i := 0; i < cap(f.children); i++ {
				size += <-f.children
			}
		}
		if f.parent != nil {
			f.parent <- size
		}
		time.Sleep(time.Millisecond)
		if f.path == root {
			println(f.path, size, f.info.Size())
		}
	}, func(filePath string, err error) {
		println(filePath, err.Error())
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWalkerWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var cnt int64
	err := Walker[int64](ctx, "..", 1, func(ctx context.Context, f *file[int64]) {
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt64(&cnt, 1)
	}, func(filePath string, err error) {
		println(filePath, err.Error())
	})
	if err == nil {
		t.Error(fmt.Errorf("expected (%s), actual (nil)", context.DeadlineExceeded.Error()))
		return
	}
	if err != context.DeadlineExceeded {
		t.Error(fmt.Errorf("expected (%s), actual (%s)", context.DeadlineExceeded.Error(), err.Error()))
		return
	}
	if atomic.LoadInt64(&cnt) != 2 {
		t.Error(fmt.Errorf("cnt should be 2, actual (%d)", cnt))
		return
	}
}
