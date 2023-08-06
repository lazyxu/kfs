package client

import (
	"fmt"
	"github.com/lazyxu/kfs/core"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"

	"github.com/muesli/termenv"
)

type TerminalUploadProcess struct {
	ch         chan *core.Process
	wg         sync.WaitGroup
	concurrent int
	conns      []net.Conn
}

func (h *TerminalUploadProcess) Show(p *core.Process) {
	h.ch <- p
}

func (h *TerminalUploadProcess) Verbose() bool {
	return true
}

func (h *TerminalUploadProcess) StartFile(index int, filePath string, info os.FileInfo) {
}

func (h *TerminalUploadProcess) OnFileError(index int, filePath string, info os.FileInfo, err error) {
	h.ch <- &core.Process{
		FilePath:  filePath,
		Err:       err,
		StackSize: -1,
	}
}

func (h *TerminalUploadProcess) StackSizeHandler(size int) {
	h.ch <- &core.Process{
		StackSize: size,
	}
}

type LineProcess struct {
	port  string
	size  uint64
	count int
}

func (h *TerminalUploadProcess) New(srcPath string, concurrent int, conns []net.Conn) core.UploadProcess {
	h.ch = make(chan *core.Process)
	h.wg.Add(1)
	h.conns = conns
	h.concurrent = concurrent
	go func() {
		h.handleProcess(srcPath)
		h.wg.Done()
	}()
	return h
}

func (h *TerminalUploadProcess) Close(resp core.FileResp, err error) {
	close(h.ch)
	h.wg.Wait()
}

func (h *TerminalUploadProcess) EnqueueFile(info os.FileInfo) {
}

func (h *TerminalUploadProcess) EndFile(index int, filePath string, info os.FileInfo, exist bool) {
}

func (h *TerminalUploadProcess) handleProcess(srcPath string) {
	lines := make([]*LineProcess, h.concurrent)
	errCnt := 0

	println()
	for i := 0; i < h.concurrent; i++ {
		println()
	}
	println()
	for p := range h.ch {
		if p == nil {
			break
		}
		rel, _ := filepath.Rel(srcPath, p.FilePath)
		if p.Err != nil {
			println(rel+":", p.Err.Error())
			errCnt++
			continue
		}
		if p.StackSize != -1 {
			size := h.concurrent
			offset := size + 2 + errCnt
			termenv.CursorPrevLine(offset)
			termenv.ClearLine()
			fmt.Printf("waiting to process: %d", p.StackSize)
			termenv.CursorNextLine(offset)
			continue
		}
		port := h.conns[p.Index].LocalAddr().String()
		port = port[strings.LastIndexByte(port, ':')+1:]

		if lines[p.Index] == nil {
			lines[p.Index] = &LineProcess{
				port: port,
			}
		}
		line := lines[p.Index]

		size := h.concurrent
		if p.Label == "code=0" || p.Label == "exist" {
			line.size += p.Size
			line.count++
		}
		offset := size + 1 - p.Index + errCnt
		termenv.CursorPrevLine(offset)
		termenv.ClearLine()
		fmt.Printf("%5s %6s %d: %-8s %6s %s", port, humanize.Bytes(line.size), line.count, p.Label, humanize.Bytes(p.Size), rel)
		termenv.CursorNextLine(offset)
	}
}
