package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lazyxu/kfs/rpc/grpcclient"

	"github.com/schollz/progressbar/v3"

	"github.com/lazyxu/kfs/core"
)

type UploadProcessBar struct {
	bar *progressbar.ProgressBar
}

func (process *UploadProcessBar) New(max int, filename string) core.UploadProcess {
	bar := progressbar.NewOptions(max,
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
		progressbar.OptionSetDescription("[1/2][hash] "+grpcclient.FormatFilename(filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return &UploadProcessBar{bar}
}

func (process *UploadProcessBar) Close() error {
	return nil
}

func (process *UploadProcessBar) BeforeContent(hash string, filename string) {
	process.bar.Reset()
	process.bar.Describe("[2/2][" + hash[0:4] + "] " + grpcclient.FormatFilename(filename))
}

func (process *UploadProcessBar) MultiWriter(w io.Writer) io.Writer {
	return io.MultiWriter(process.bar, w)
}
