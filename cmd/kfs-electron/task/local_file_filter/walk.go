package local_file_filter

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"os"
	"path/filepath"
	"sync"
)

type DriverLocalFile struct {
	driverId uint64
	mutex    sync.Locker
	taskInfo TaskInfo
}

type filterHandlers struct {
	core.DefaultWalkDirHandlers
	d        *DriverLocalFile
	driverId uint64
	srcPath  string
}

func (h *filterHandlers) FilePathFilter(filePath string) bool {
	return false
}

func (h *filterHandlers) OnFileError(filePath string, err error) {
	h.d.addTaskWarning(err.Error())
	println(filePath+":", err.Error())
}

func (h *filterHandlers) DirHandler(ctx context.Context, filePath string, dirInfo os.FileInfo, infos []os.FileInfo, continues []bool) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	if filePath != h.srcPath {
		h.d.setTaskDir(filePath, uint64(len(infos)))
	}

	for _, info := range infos {
		h.d.addTaskTotal(info)
	}

	if filePath != h.srcPath {
		h.d.addTaskCnt(dirInfo)
	}

	return nil
}

func (d *DriverLocalFile) checkFilter(ctx context.Context, driverId uint64, srcPath string) error {
	if !filepath.IsAbs(srcPath) {
		return errors.New("请输入绝对路径")
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("源目录不存在")
	}
	fmt.Println("filter start", driverId, srcPath)

	handlers := &filterHandlers{
		d:        d,
		driverId: driverId,
		srcPath:  srcPath,
	}
	err = core.WalkDir(ctx, srcPath, handlers)
	if err != nil {
		return err
	}

	fmt.Printf("filter finish %+v\n", d)
	return nil
}
