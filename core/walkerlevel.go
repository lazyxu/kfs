package core

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/queues/arrayqueue"
	"os"
	"path/filepath"
	"sync"
)

type file struct {
	Path     string
	Info     os.FileInfo
	parent   chan error
	children []*file
}

type WalkByLevelHandlers interface {
	FilePathFilter(filePath string) bool
	FileInfoFilter(filePath string, info os.FileInfo) bool
	FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo) error
	OnFileError(filePath string, index int, info os.FileInfo, err error)
	StartWorker(ctx context.Context, index int)
	EndWorker(ctx context.Context, index int)
	AddToWorkList(info os.FileInfo)
	HasEnqueuedAll()
}

type DefaultWalkByLevelHandlers struct{}

func (DefaultWalkByLevelHandlers) FilePathFilter(filePath string) bool {
	return false
}

func (DefaultWalkByLevelHandlers) FileInfoFilter(filePath string, info os.FileInfo) bool {
	return false
}

func (DefaultWalkByLevelHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo) error {
	return nil
}

func (DefaultWalkByLevelHandlers) OnFileError(filePath string, info os.FileInfo, err error) {
	println(filePath, err.Error())
}

func (DefaultWalkByLevelHandlers) StartWorker(ctx context.Context, index int) {
}

func (DefaultWalkByLevelHandlers) EndWorker(ctx context.Context, index int) {
}

func (DefaultWalkByLevelHandlers) AddToWorkList(info os.FileInfo) {
}
func (DefaultWalkByLevelHandlers) HasEnqueuedAll() {
}

func WalkByLevel(ctx context.Context, filePath string, concurrent int, handlers WalkByLevelHandlers) (err1 error) {
	if concurrent <= 0 {
		return fmt.Errorf("concurrent should be > 0, actual %d", concurrent)
	}
	filePath, err1 = filepath.Abs(filePath)
	if err1 != nil {
		return
	}
	root := &file{
		Path:   filePath,
		parent: make(chan error, 1),
	}
	root.parent <- nil
	queue := arrayqueue.New()
	queue.Enqueue(root)
	ch := make(chan *file, 1000000)
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
				select {
				case err := <-f.parent:
					if err != nil {
						continue
					}
				case <-ctx.Done():
					err1 = context.DeadlineExceeded
					goto EndWorker
				}
				err := handlers.FileHandler(ctx, index, f.Path, f.Info)
				if err != nil {
					handlers.OnFileError(f.Path, index, f.Info, err)
				}
				for _, child := range f.children {
					select {
					case child.parent <- err:
					case <-ctx.Done():
						err1 = context.DeadlineExceeded
						goto EndWorker
					}
				}
			}
		EndWorker:
			handlers.EndWorker(ctx, index)
			wg.Done()
		}(i)
	}
	for !queue.Empty() {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		v, _ := queue.Dequeue()
		f := v.(*file)
		var info os.FileInfo
		info, continues := filters(f.Path, handlers)
		if continues {
			continue
		}
		f.Info = info

		if info.IsDir() {
			infos, err := os.ReadDir(f.Path)
			if err != nil {
				handlers.OnFileError(f.Path, -1, f.Info, err)
				continue
			}
			for i := len(infos) - 1; i >= 0; i-- {
				child := &file{
					Path:   filepath.Join(f.Path, infos[i].Name()),
					parent: make(chan error, 1),
				}
				queue.Enqueue(child)
				f.children = append(f.children, child)
			}
		}
		handlers.AddToWorkList(f.Info)
		ch <- f
	}
	handlers.HasEnqueuedAll()
	return
}

func filters(filePath string, handlers WalkByLevelHandlers) (os.FileInfo, bool) {
	if handlers.FilePathFilter(filePath) {
		return nil, true
	}
	info, err := os.Lstat(filePath)
	if err != nil {
		handlers.OnFileError(filePath, -1, info, err)
		return nil, true
	}
	if handlers.FileInfoFilter(filePath, info) {
		return info, true
	}
	return info, false
}
