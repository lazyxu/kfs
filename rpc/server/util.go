package server

import (
	"archive/zip"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
)

func getHash(f *zip.File) (string, string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", "", err
	}
	defer rc.Close()

	// 创建两个哈希计算器
	hashSHA256 := sha256.New()
	hashMD5 := md5.New()

	// 创建管道：将数据同时分流到 SHA256 和 MD5
	pipeReader, pipeWriter := io.Pipe()
	tee := io.TeeReader(rc, pipeWriter) // 分流到管道写入端

	// 启动协程：从管道读取端计算 MD5
	var wg sync.WaitGroup
	wg.Add(1)
	var md5Sum string
	var md5Err error
	go func() {
		defer wg.Done()
		_, md5Err = io.Copy(hashMD5, pipeReader)
		md5Sum = hex.EncodeToString(hashMD5.Sum(nil))
	}()

	// 主线程：从 tee 读取并计算 SHA256
	_, err = io.Copy(hashSHA256, tee) // 读取时同步写入管道
	pipeWriter.Close()                // 确保管道写入结束
	if err != nil {
		return "", "", err
	}

	// 等待 MD5 计算完成
	wg.Wait()
	if md5Err != nil {
		return "", "", md5Err
	}

	shaSum := hex.EncodeToString(hashSHA256.Sum(nil))
	return md5Sum, shaSum, nil
}

func Orientation(img image.Image, o int64) image.Image {
	switch o {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.Rotate180(imaging.FlipH(img))
	case 5:
		return imaging.Rotate90(imaging.FlipH(img))
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Rotate270(imaging.FlipH(img))
	case 8:
		return imaging.Rotate90(img)
	}
	return img
}

var (
	ffmpegMutex sync.Mutex
)

func SafeFFmpegRun(ctx context.Context, src, dst string) error {
	ffmpegMutex.Lock()
	defer ffmpegMutex.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	tmpPath := filepath.Join(os.TempDir(), "ffmpeg.mp4")
	defer os.Remove(tmpPath)

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	numCPU := runtime.NumCPU()
	fmt.Printf("Running ffmpeg command: ffmpeg -noautorotate -threads %d -i %s -qscale 0 -vcodec libx264 %s\n", numCPU, src, tmpPath)
	cmd := exec.CommandContext(cancelCtx, "ffmpeg",
		"-noautorotate", "-threads", strconv.Itoa(numCPU), "-i", src,
		"-qscale", "0", "-vcodec", "libx264", tmpPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Rename(tmpPath, dst)
}
