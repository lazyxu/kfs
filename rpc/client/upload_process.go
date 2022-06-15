package client

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/muesli/termenv"
)

type Process struct {
	index     int
	label     string
	filePath  string
	size      uint64
	stackSize int
	err       error
}

func (h *uploadHandlers) ErrHandler(filePath string, err error) {
	if h.verbose {
		h.ch <- &Process{
			filePath:  filePath,
			err:       err,
			stackSize: -1,
		}
	} else {
		println(filePath+":", err.Error())
	}
}

func (h *uploadHandlers) StackSizeHandler(size int) {
	if h.verbose {
		h.ch <- &Process{
			stackSize: size,
		}
	}
}

type LineProcess struct {
	port  string
	size  uint64
	count int
}

func (h *uploadHandlers) handleProcess(srcPath string) {
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
		rel, _ := filepath.Rel(srcPath, p.filePath)
		if p.err != nil {
			println(rel+":", p.err.Error())
			errCnt++
			continue
		}
		if p.stackSize != -1 {
			size := h.concurrent
			offset := size + 2 + errCnt
			termenv.CursorPrevLine(offset)
			termenv.ClearLine()
			fmt.Printf("waiting to process: %d", p.stackSize)
			termenv.CursorNextLine(offset)
			continue
		}
		port := h.conns[p.index].LocalAddr().String()
		port = port[strings.LastIndexByte(port, ':')+1:]

		if lines[p.index] == nil {
			lines[p.index] = &LineProcess{
				port: port,
			}
		}
		line := lines[p.index]

		size := h.concurrent
		if p.label == "code=0" || p.label == "exist" {
			line.size += p.size
			line.count++
		}
		offset := size + 1 - p.index + errCnt
		termenv.CursorPrevLine(offset)
		termenv.ClearLine()
		fmt.Printf("%5s %6s %d: %-8s %6s %s", port, humanize.Bytes(line.size), line.count, p.label, humanize.Bytes(p.size), rel)
		termenv.CursorNextLine(offset)
	}
}
