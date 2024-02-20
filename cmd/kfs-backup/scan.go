package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/lazyxu/kfs/core"
	ignore "github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
	"time"
)

type scanHandlers struct {
	core.DefaultWalkDirHandlers
	srcPath     string
	gitIgnore   *ignore.GitIgnore
	ignoreCount int
	errCount    int
	size        uint64
	fileCount   uint64
	dirCount    uint64
	verbose     bool
}

func (h *scanHandlers) FilePathFilter(filePath string) bool {
	ignored := h.gitIgnore.MatchesPath(filePath)
	if ignored {
		h.ignoreCount++
		if h.verbose {
			rel, _ := filepath.Rel(h.srcPath, filePath)
			fmt.Printf("跳过第 %d 个文件或目录 - %s\n", h.ignoreCount, rel)
		}
	}
	return ignored
}

func (h *scanHandlers) OnFileError(filePath string, err error) {
	h.errCount++
	fmt.Printf("发现第 %d 个错误 - %s: %s\n", h.errCount, filePath, err.Error())
}

func (h *scanHandlers) DirHandler(ctx context.Context, filePath string, dirInfo os.FileInfo, infos []os.FileInfo, continues []bool) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	for _, info := range infos {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
		if info.IsDir() {
			h.dirCount++
			if h.verbose {
				rel, _ := filepath.Rel(h.srcPath, filepath.Join(filePath, info.Name()))
				fmt.Printf("扫描到第 %d 个目录 - %s\n", h.dirCount, rel)
			}
		} else {
			h.fileCount++
			size := uint64(info.Size())
			h.size += size
			if h.verbose {
				rel, _ := filepath.Rel(h.srcPath, filepath.Join(filePath, info.Name()))
				fmt.Printf("扫描到第 %d 个文件，大小为 %s - %s\n", h.fileCount, humanize.IBytes(size), rel)
			}
		}
	}

	return nil
}

func doScan(ctx context.Context, srcPath string, ignores []string, verbose bool) {
	start := time.Now()
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		fmt.Printf("路径格式错误：%s\n", srcPath)
		return
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		fmt.Printf("查看目录状态失败：%s\n", err.Error())
		return
	}
	if !info.IsDir() {
		fmt.Printf("请输入一个目录：%s\n", srcPath)
		return
	}
	fmt.Printf("开始扫描：%s\n", srcPath)
	//if os.PathSeparator == '\\' {
	//	ignores = strings.ReplaceAll(ignores, "\\", "/")
	//}
	//list := strings.Split(ignores, "\n")
	gitIgnore := ignore.CompileIgnoreLines(ignores...)
	handlers := &scanHandlers{
		srcPath:   srcPath,
		gitIgnore: gitIgnore,
		verbose:   verbose,
	}
	err = core.WalkDir(ctx, srcPath, handlers)
	if err != nil {
		fmt.Printf("扫描时发生错误：%s\n", err.Error())
		return
	}

	fmt.Printf("扫描完成！耗时 %s，跳过了 %d 个文件或目录，", time.Since(start).String(), handlers.ignoreCount)
	if handlers.errCount > 0 {
		fmt.Printf("扫描时发生 %d 处错误\n", handlers.errCount)
	} else {
		fmt.Printf("扫描时未发生错误\n")
	}
	fmt.Printf("总大小：%s，共 %d 个文件，%d 个目录\n", humanize.IBytes(handlers.size), handlers.fileCount, handlers.dirCount)
	return
}
