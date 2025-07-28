package server

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
)

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
