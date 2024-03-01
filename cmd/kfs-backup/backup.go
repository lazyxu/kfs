package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	ignore "github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
	"time"
)

type WebUploadDirProcess struct {
	srcPath      string
	gitIgnore    *ignore.GitIgnore
	ignoreCount  int
	errCount     int
	newSize      uint64
	newFileCount uint64
	size         uint64
	count        uint64
	dirCount     uint64
	verbose      bool
}

var _ core.UploadDirProcess = &WebUploadDirProcess{}

func (h *WebUploadDirProcess) FilePathFilter(filePath string) bool {
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

func (h *WebUploadDirProcess) OnFileError(filePath string, err error) {
	h.errCount++
	fmt.Printf("发现第 %d 个错误 - %s: %s\n", h.errCount, filePath, err.Error())
}

func (h *WebUploadDirProcess) StartFile(filePath string, info os.FileInfo) {
	//if h.verbose {
	//	size := uint64(info.Size())
	//	rel, _ := filepath.Rel(h.srcPath, filePath)
	//	fmt.Printf("开始上传文件，大小为 %s - %s\n", humanize.IBytes(size), rel)
	//}
}

func (h *WebUploadDirProcess) StartUploadFile(filePath string, info os.FileInfo, hash string) {
	if h.verbose {
		size := uint64(info.Size())
		rel, _ := filepath.Rel(h.srcPath, filePath)
		fmt.Printf("第 %5d 个文件 %s 开始上传，大小为 %s, 哈希值为%s\n", h.count, rel, humanize.IBytes(size), hash)
	}
}

func (h *WebUploadDirProcess) EndUploadFile(filePath string, info os.FileInfo) {
	size := uint64(info.Size())
	h.count++
	h.size += size
	h.newFileCount++
	h.newSize += size
	if h.verbose {
		rel, _ := filepath.Rel(h.srcPath, filePath)
		fmt.Printf("第 %5d 个文件 %s 上传完成\n", h.count, rel)
	}
}

func (h *WebUploadDirProcess) SkipFile(filePath string, info os.FileInfo, hash string) {
	size := uint64(info.Size())
	h.count++
	h.size += size
	if h.verbose {
		rel, _ := filepath.Rel(h.srcPath, filePath)
		fmt.Printf("第 %5d 个文件 %s 已存在，跳过, 哈希值为 %s\n", h.count, rel, hash)
	}
}

func (h *WebUploadDirProcess) EndFile(filePath string, info os.FileInfo) {
	//h.fileCount++
	//size := uint64(info.Size())
	//h.size += size
	//if h.verbose {
	//	rel, _ := filepath.Rel(h.srcPath, filePath)
	//	fmt.Printf("第 %d 个文件上传完成 - %s\n", h.fileCount, rel)
	//}
}

func (h *WebUploadDirProcess) StartDir(filePath string, n uint64) {
	h.dirCount++
	if h.verbose {
		rel, _ := filepath.Rel(h.srcPath, filePath)
		fmt.Printf("第 %5d 个目录 %s 开始上传, 目录下包括 %d 个文件或子目录\n", h.dirCount, rel, n)
	}
}

func (h *WebUploadDirProcess) EndDir(filePath string, info os.FileInfo) {
	if h.verbose {
		rel, _ := filepath.Rel(h.srcPath, filePath)
		fmt.Printf("第 %5d 个目录 %s 上传完成\n", h.dirCount, rel)
	}
}

func (h *WebUploadDirProcess) PushFile(info os.FileInfo) {
}

func doUpload(ctx context.Context, deviceId string, serverAddr string, driverId uint64, srcPath string, ignores []string, verbose bool) {
	var encoder string
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
	fmt.Printf("开始上传 %s 到 %s 的 云盘%d 上\n", srcPath, serverAddr, driverId)

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}
	//if os.PathSeparator == '\\' {
	//	ignores = strings.ReplaceAll(ignores, "\\", "/")
	//}
	//list := strings.Split(ignores, "\n")
	gitIgnore := ignore.CompileIgnoreLines(ignores...)
	handlers := &WebUploadDirProcess{
		srcPath:   srcPath,
		gitIgnore: gitIgnore,
		verbose:   verbose,
	}
	err = fs.UploadDir(ctx, deviceId, driverId, "/", srcPath, core.UploadDirConfig{
		UploadDirProcess: handlers,
		Encoder:          encoder,
		Concurrent:       1,
		Verbose:          verbose,
	})
	if err != nil {
		fmt.Printf("上传时发生错误：%s\n", err.Error())
		return
	}

	fmt.Printf("上传完成！耗时 %s，忽略了 %d 个文件或目录，", time.Since(start).String(), handlers.ignoreCount)
	if handlers.errCount > 0 {
		fmt.Printf("上传时发生 %d 处错误\n", handlers.errCount)
	} else {
		fmt.Printf("上传时未发生错误\n")
	}
	fmt.Printf("本次上传 %s, 包括 %d 个文件\n", humanize.IBytes(handlers.newSize), handlers.newFileCount)
	fmt.Printf("总大小 %s，包括 %d 个文件，%d 个目录\n", humanize.IBytes(handlers.size), handlers.count, handlers.dirCount)
	return
}
