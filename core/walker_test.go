package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

type TestWalkerHandlers struct {
	DefaultWalkHandlers[int64]
}

func (TestWalkerHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []int64) int64 {
	var size int64
	if !info.IsDir() {
		size = info.Size()
	}
	for _, child := range children {
		size += child
	}
	return size
}

func TestWalker(t *testing.T) {
	ctx := context.Background()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := Walk[int64](ctx, root, 15, &TestWalkerHandlers{})
	if err != nil {
		t.Error(err)
		return
	}
	println(root, resp)
}

type TestWalkerWithTimeoutHandlers struct {
	DefaultWalkHandlers[int64]
	cnt int64
}

func (h *TestWalkerWithTimeoutHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []int64) int64 {
	time.Sleep(100 * time.Millisecond)
	atomic.AddInt64(&h.cnt, 1)
	return 0
}

func TestWalkerWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	handlers := &TestWalkerWithTimeoutHandlers{}
	_, err := Walk[int64](ctx, ".", 1, handlers)
	if err == nil {
		t.Error(fmt.Errorf("expected (%s), actual (nil)", context.DeadlineExceeded.Error()))
		return
	}
	if err != context.DeadlineExceeded {
		t.Error(fmt.Errorf("expected (%s), actual (%s)", context.DeadlineExceeded.Error(), err.Error()))
		return
	}
	if atomic.LoadInt64(&handlers.cnt) != 2 {
		t.Error(fmt.Errorf("cnt should be 2, actual (%d)", atomic.LoadInt64(&handlers.cnt)))
		return
	}
}

type TestWalkerPathFilterHandlers struct {
	DefaultWalkHandlers[int64]
	root string
}

func (h *TestWalkerPathFilterHandlers) FilePathFilter(filePath string) bool {
	return filePath == h.root
}

func TestWalkerPathFilter(t *testing.T) {
	ctx := context.Background()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := Walk[int64](ctx, root, 15, &TestWalkerPathFilterHandlers{root: root})
	if err != nil {
		t.Error(err)
		return
	}
	if resp != 0 {
		t.Error(fmt.Errorf("filter root: expected 0, actual (%d)", resp))
		return
	}
}
