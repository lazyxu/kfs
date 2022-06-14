package client

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/muesli/termenv"

	"github.com/emirpasic/gods/sets/linkedhashset"
)

type Process struct {
	label     string
	conn      net.Conn
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
	index int
	port  string
	size  uint64
	count int
}

func addToSet(set *linkedhashset.Set, port string) (int, *LineProcess) {
	index, lp := set.Find(func(index int, value interface{}) bool {
		v := value.(*LineProcess)
		return v.port == port
	})
	if lp != nil {
		return index, lp.(*LineProcess)
	}
	line := &LineProcess{
		index: set.Size(),
		port:  port,
	}
	set.Add(line)
	return index, line
}

func (h *uploadHandlers) handleProcess(srcPath string, concurrent int) {
	set := linkedhashset.New()
	errCnt := 0

	println()
	for i := 0; i < concurrent; i++ {
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
			size := set.Size()
			offset := size + 2 + errCnt
			termenv.CursorPrevLine(offset)
			termenv.ClearLine()
			fmt.Printf("waiting to process: %d", p.stackSize)
			termenv.CursorNextLine(offset)
			continue
		}
		port := p.conn.LocalAddr().String()
		port = port[strings.LastIndexByte(port, ':')+1:]
		_, line := addToSet(set, port)
		size := set.Size()
		if p.label == "code=0" || p.label == "exist" {
			line.size += p.size
			line.count++
		}
		offset := size + 1 - line.index + errCnt
		termenv.CursorPrevLine(offset)
		termenv.ClearLine()
		fmt.Printf("%5s %6s %d: %-8s %6s %s", port, humanize.Bytes(line.size), line.count, p.label, humanize.Bytes(p.size), rel)
		termenv.CursorNextLine(offset)
	}
}
