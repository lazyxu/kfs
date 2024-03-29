package local_file_filter

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	ignore "github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DriverLocalFile struct {
	driverId uint64
	mutex    sync.Locker
	taskInfo TaskInfo
}

type filterHandlers struct {
	core.DefaultWalkDirHandlers
	d         *DriverLocalFile
	driverId  uint64
	srcPath   string
	gitIgnore *ignore.GitIgnore
}

func (h *filterHandlers) FilePathFilter(filePath string) bool {
	ignored := h.gitIgnore.MatchesPath(filePath)
	if ignored {
		println(filePath + ": ignored")
		h.d.addTaskIgnores(filePath)
	} else {
		println(filePath)
	}
	return ignored
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
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
		h.d.addTaskTotal(info)
	}

	if filePath != h.srcPath {
		h.d.addTaskCnt(dirInfo)
	}

	return nil
}

func (d *DriverLocalFile) DoFilter(ctx context.Context, srcPath, ignores string) error {
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
	fmt.Println("filter start", d.driverId, srcPath, ignores)
	if os.PathSeparator == '\\' {
		ignores = strings.ReplaceAll(ignores, "\\", "/")
	}
	list := strings.Split(ignores, "\n")
	gitIgnore := ignore.CompileIgnoreLines(list...)
	handlers := &filterHandlers{
		d:         d,
		driverId:  d.driverId,
		srcPath:   srcPath,
		gitIgnore: gitIgnore,
	}
	err = core.WalkDir(ctx, srcPath, handlers)
	if err != nil {
		return err
	}

	fmt.Printf("filter finish %+v\n", d)
	return nil
}
