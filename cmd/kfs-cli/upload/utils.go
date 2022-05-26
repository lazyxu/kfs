package upload

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/lazyxu/kfs/pb"
)

func SendHeader(hash string, name string, info os.FileInfo, fn func(item *pb.DirItem) error) error {
	modifyTime := uint64(info.ModTime().UnixNano())
	err := fn(&pb.DirItem{
		Hash:       hash,
		Name:       name,
		Mode:       uint64(info.Mode()),
		Size:       uint64(info.Size()),
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	})
	if err != nil {
		return err
	}
	return err
}

func SendContent(bar *progressbar.ProgressBar, hash string, filename string, fn func(data []byte, isFirst bool, isLast bool) error) error {
	bar.Reset()
	bar.Describe("[2/2][" + hash[0:4] + "] " + FormatFilename(filename))
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	isFirst := true
	for {
		chunk := make([]byte, 0, fileChunkSize)
		var n int64
		w := io.MultiWriter(bytes.NewBuffer(chunk), bar)
		n, err = io.Copy(w, io.LimitReader(f, fileChunkSize))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = fn(chunk[:n], isFirst, n < fileChunkSize)
		isFirst = false
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n < fileChunkSize {
			break
		}
	}
	return nil
}

func GetFileHash(bar io.Writer, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := sha256.New()
	w := io.MultiWriter(hash, bar)
	_, err = io.Copy(w, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func NewProcessBar(info os.FileInfo, filename string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(int(info.Size()),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(20*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetDescription("[1/2][hash] "+FormatFilename(filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}

func FormatFilename(filename string) string {
	var name = []rune(filepath.Base(filename))
	if len(name) > 10 {
		name = append(name[:10], []rune("..")...)
	}
	return fmt.Sprintf("%-12s", string(name))
}
