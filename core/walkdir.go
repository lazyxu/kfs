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
	DirHandler(ctx context.Context, filePath string, infos []os.FileInfo, continues []bool) error
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

func (DefaultWalkDirHandlers) DirHandler(ctx context.Context, filePath string, infos []os.FileInfo, continues []bool) error {
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
	return handleDir(ctx, filePath, handlers)
}

func handleDir(ctx context.Context, filePath string, handlers WalkDirHandlers) error {
	infos, err := os.ReadDir(filePath)
	if err != nil {
		handlers.OnFileError(filePath, err)
		return err
	}
	filteredInfos := make([]os.FileInfo, len(infos))
	for i, info := range infos {
		fp := filepath.Join(filePath, info.Name())
		if handlers.FilePathFilter(fp) {
			continue
		}
		info2, err1 := info.Info()
		if err1 != nil {
			handlers.OnFileError(fp, err1)
			continue
		}
		if handlers.FileInfoFilter(fp, info2) {
			continue
		}
		filteredInfos[i] = info2
	}
	continues := make([]bool, len(filteredInfos))
	err = handlers.DirHandler(ctx, filePath, filteredInfos, continues)
	if err != nil {
		handlers.OnFileError(filePath, err)
	}
	for i, fi := range filteredInfos {
		if continues[i] || !fi.IsDir() {
			continue
		}
		fp := filepath.Join(filePath, fi.Name())
		err = handleDir(ctx, fp, handlers)
		if err != nil {
			handlers.OnFileError(filePath, err)
		}
	}
	return nil
}
