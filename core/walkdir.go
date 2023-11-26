package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

type WalkDirHandlers interface {
	FilePathFilter(filePath string) bool
	FileInfoFilter(filePath string, info os.FileInfo) bool
	OnFileError(filePath string, err error)
	DirHandler(ctx context.Context, filePath string, dirInfo os.FileInfo, infos []os.FileInfo, continues []bool) error
}

type DefaultWalkDirHandlers struct{}

var _ WalkDirHandlers = DefaultWalkDirHandlers{}

func (DefaultWalkDirHandlers) FilePathFilter(filePath string) bool {
	return false
}

func (DefaultWalkDirHandlers) FileInfoFilter(filePath string, info os.FileInfo) bool {
	return false
}

func (DefaultWalkDirHandlers) OnFileError(filePath string, err error) {
	println(filePath, err.Error())
}

func (DefaultWalkDirHandlers) DirHandler(ctx context.Context, filePath string, dirInfo os.FileInfo, infos []os.FileInfo, continues []bool) error {
	return nil
}

func WalkDir(ctx context.Context, filePath string, handlers WalkDirHandlers) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	if handlers.FilePathFilter(filePath) {
		return nil
	}
	info, err := os.Lstat(filePath)
	if err != nil {
		handlers.OnFileError(filePath, err)
		return err
	}
	if handlers.FileInfoFilter(filePath, info) {
		return nil
	}
	if !info.IsDir() {
		return errors.New("expected dir path")
	}
	return handleDir(ctx, filePath, info, handlers)
}

func handleDir(ctx context.Context, filePath string, dirInfo os.FileInfo, handlers WalkDirHandlers) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	dirEntries, err := os.ReadDir(filePath)
	if err != nil {
		handlers.OnFileError(filePath, err)
		return err
	}
	filteredInfos := []os.FileInfo{}
	for _, dirEntry := range dirEntries {
		fp := filepath.Join(filePath, dirEntry.Name())
		if handlers.FilePathFilter(fp) {
			continue
		}
		info, err1 := dirEntry.Info()
		if err1 != nil {
			handlers.OnFileError(fp, err1)
			continue
		}
		if handlers.FileInfoFilter(fp, info) {
			continue
		}
		filteredInfos = append(filteredInfos, info)
	}
	continues := make([]bool, len(filteredInfos))
	err = handlers.DirHandler(ctx, filePath, dirInfo, filteredInfos, continues)
	if errors.Is(err, context.Canceled) {
		return err
	}
	if err != nil {
		handlers.OnFileError(filePath, err)
	}
	for i, fi := range filteredInfos {
		if continues[i] || !fi.IsDir() {
			continue
		}
		fp := filepath.Join(filePath, fi.Name())
		err = handleDir(ctx, fp, fi, handlers)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			handlers.OnFileError(filePath, err)
		}
	}
	return nil
}
