package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/dbBase"
)

type CountAndSize struct {
	Count int64
	Size  int64
}

type FileSizeResp struct {
	FileSize  int64 `json:"fileSize"`
	FileCount int64 `json:"fileCount"`
	DirCount  int64 `json:"dirCount"`
	StackSize int   `json:"stackSize"`
}

type SizeWalkerHandlers struct {
	FileSizeResp
	core.DefaultWalkHandlers[CountAndSize]
	req         WsReq
	onResp      func(finished bool, data interface{}) error
	tick        <-chan time.Time
	db          *DB
	lock        sync.Locker
	dbFileInfos []DbFileInfo
}

func (h *SizeWalkerHandlers) StackSizeHandler(size int) {
	h.StackSize = size
}

type DbFileInfo struct {
	path  string
	typ   int // 0: file 1: dir
	count int64
	size  int64
}

func (h *SizeWalkerHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []CountAndSize) CountAndSize {
	var count int64 = 1
	var size int64
	if info.IsDir() {
		atomic.AddInt64(&h.DirCount, 1)
		for _, child := range children {
			count += child.Count
			size += child.Size
		}
	} else {
		count = 1
		size = info.Size()
		atomic.AddInt64(&h.FileCount, 1)
		atomic.AddInt64(&h.FileSize, info.Size())
	}

	h.addFile(info, filePath, count, size)

	//err := h.db.InsertFile(ctx, h.startTime, filePath, info.IsDir(), count, size)
	//if err != nil {
	//	panic(err)
	//}

	select {
	case <-h.tick:
		fmt.Printf("tick: %+v\n", h.FileSizeResp)
		err := h.onResp(false, h.FileSizeResp)
		if err != nil {
			fmt.Printf("%+v %+v\n", h.req, err)
		}
	case <-ctx.Done():
	default:
	}
	return CountAndSize{
		Count: count,
		Size:  size,
	}
}

func (h *SizeWalkerHandlers) addFile(info os.FileInfo, filePath string, count int64, size int64) {
	{
		typ := 0
		if info.IsDir() {
			typ = 1
		}
		h.lock.Lock()
		defer h.lock.Unlock()
		h.dbFileInfos = append(h.dbFileInfos, DbFileInfo{
			path:  filePath,
			typ:   typ,
			count: count,
			size:  size,
		})
	}
}

func (h *SizeWalkerHandlers) insertFiles(ctx context.Context, id int64) error {
	conn := h.db.getConn()
	defer h.db.putConn(conn)
	return dbBase.InsertBatch[DbFileInfo](ctx, conn, 32766, h.dbFileInfos, 7, getInsertItemQuery, func(args []interface{}, start int, item DbFileInfo) {
		args[start] = id
		args[start+1] = item.path
		args[start+2] = filepath.Dir(item.path)
		args[start+3] = filepath.Base(item.path)
		args[start+4] = item.typ
		args[start+5] = item.count
		args[start+6] = item.size
	})
}

func (h *SizeWalkerHandlers) insertScanHistory(ctx context.Context, startTime int64, dirname string, fileSize int64, fileCount int64, dirCount int64) (int64, error) {
	conn := h.db.getConn()
	defer h.db.putConn(conn)
	res, err := conn.ExecContext(ctx, `
	INSERT INTO _scan_history (
	    time,
	    dirname,
	    fileSize,
	    fileCount,
	    dirCount
	) VALUES (?, ?, ?, ?, ?);`, startTime, dirname, fileSize, fileCount, dirCount)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getInsertItemQuery(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
	INSERT INTO _file (
	    id,
		path,
	    dirname,
		name,
	    typ,
		count,
		size
	) VALUES `)
	if err != nil {
		return "", err
	}
	for i := 0; i < row; i++ {
		if i != 0 {
			qs.WriteString(", ")
		}
		qs.WriteString("(?, ?, ?, ?, ?, ?, ?)")
	}
	qs.WriteString(";")
	return qs.String(), err
}

func (p *WsProcessor) scan(ctx context.Context, db *DB, req WsReq, srcPath string, concurrent int) error {
	if !filepath.IsAbs(srcPath) {
		return p.err(req, errors.New("请输入绝对路径"))
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		return p.err(req, err)
	}
	if !info.IsDir() {
		return p.err(req, errors.New("请输入一个目录"))
	}
	startTime := time.Now().UnixNano()
	w := SizeWalkerHandlers{
		req:  req,
		tick: time.Tick(time.Millisecond * 500),
		onResp: func(finished bool, data interface{}) error {
			return p.ok(req, finished, data)
		},
		db:   db,
		lock: &sync.Mutex{},
	}
	err = p.ok(req, false, w.FileSizeResp)
	if err != nil {
		return err
	}
	_, err = core.Walk[CountAndSize](ctx, srcPath, concurrent, &w)
	if err != nil {
		return p.err(req, err)
	}
	err = p.ok(req, false, w.FileSizeResp)
	if err != nil {
		return err
	}
	id, err := w.insertScanHistory(ctx, startTime, srcPath, w.FileSize, w.FileCount, w.DirCount)
	if err != nil {
		return p.err(req, err)
	}
	err = w.insertFiles(ctx, id)
	if err != nil {
		return p.err(req, err)
	}
	return p.ok(req, true, w.FileSizeResp)
}
