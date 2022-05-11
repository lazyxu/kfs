package local

import (
	"context"
	"os"
)

type Visitor[T any] interface {
	Enter(filename string, info os.FileInfo) bool
	HasExit() bool
	Exit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []T) (T, error)
}

type EmptyVisitor[T any] struct {
}

func (v *EmptyVisitor[T]) Enter(filename string, info os.FileInfo) bool {
	return true
}

func (v *EmptyVisitor[T]) HasExit() bool {
	return false
}

func (v *EmptyVisitor[T]) Exit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []T) (T, error) {
	var t T
	return t, nil
}

type FileSizeVisitor struct {
	EmptyVisitor[any]
	MaxFileSize  int64
	IgnoredCount uint64
}

func (v *FileSizeVisitor) Enter(filename string, info os.FileInfo) bool {
	if info.Mode().IsRegular() {
		if info.Size() > v.MaxFileSize {
			v.IgnoredCount++
		}
	}
	return true
}

type CountVisitor struct {
	EmptyVisitor[any]
	File     uint64
	Dir      uint64
	Symlink  uint64
	FileSize uint64
}

func (v *CountVisitor) Enter(filename string, info os.FileInfo) bool {
	if info.IsDir() {
		v.Dir++
	} else if info.Mode().IsRegular() {
		v.File++
		v.FileSize += uint64(info.Size())
	} else if info.Mode()&os.ModeSymlink != 0 {
		v.Symlink++
	}
	return true
}
