package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/db/dbBase"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
)

func main() {
	err := rootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const PortStr = "port"

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "kfs-electron",
		Args: cobra.RangeArgs(1, 1),
		Run:  runRoot,
	}
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	cmd.PersistentFlags().StringP(PortStr, "p", "0", "local web server port")
	return cmd
}

func runRoot(cmd *cobra.Command, args []string) {
	serverAddr := args[0]

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, serverAddr)
	})

	portStr := cmd.Flag(PortStr).Value.String()
	lis, err := net.Listen("tcp", "0.0.0.0:"+portStr)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(os.Stdout, "Websocket server listening at:", lis.Addr().String())
	if err != nil {
		panic(err)
	}
	if err != nil {
		return
	}
	if os.Getenv("KFS_CONFIG_PATH") != "" {
		filePath := os.Getenv("KFS_CONFIG_PATH")
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			panic(err)
		}
		m["port"] = lis.Addr().(*net.TCPAddr).Port
		data, err = json.Marshal(m)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(filePath, data, 0o600)
		if err != nil {
			panic(err)
		}
	}
	err = http.Serve(lis, nil)
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func wsHandler(w http.ResponseWriter, r *http.Request, serverAddr string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	process(r.Context(), c)
}

type WsReq struct {
	Type string      `json:"type"`
	Id   int         `json:"id"`
	Data interface{} `json:"data"`
}

type WsResp struct {
	Id       int         `json:"id"`
	Finished bool        `json:"finished"`
	Data     interface{} `json:"data"`
	ErrMsg   string      `json:"errMsg"`
}

func process(ctx context.Context, conn *websocket.Conn) {
	println(conn.RemoteAddr().String(), "Process")
	var calculateBackupSizeCancelFunctions sync.Map
	var lock sync.Mutex
	//defer func() {
	//	if err := recover(); err != nil {
	//		println("recover", err)
	//		conn.Close()
	//	}
	//}()
	db, err := NewDb("electron.db")
	if err != nil {
		panic(err)
	}

	for {
		println(conn.RemoteAddr().String(), "ReadJSON")
		var req WsReq
		err := conn.ReadJSON(&req)
		if err == io.EOF || websocket.IsUnexpectedCloseError(err) {
			conn.Close()
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", req)
		switch req.Type {
		case "calculateBackupSize.cancel":
			cancelFunc, ok := calculateBackupSizeCancelFunctions.Load(req.Id)
			if !ok {
				continue
			}
			cancelFunc.(context.CancelFunc)()
		case "calculateBackupSize":
			newCtx, cancelFunc := context.WithCancel(ctx)
			calculateBackupSizeCancelFunctions.Store(req.Id, cancelFunc)
			go func() {
				err := calculateBackupSize(newCtx, db, req, req.Data.(map[string]interface{})["backupDir"].(string), func(finished bool, data interface{}) error {
					var resp WsResp
					resp.Id = req.Id
					resp.Finished = finished
					resp.Data = data
					lock.Lock()
					defer lock.Unlock()
					return conn.WriteJSON(resp)
				}, func(err error) error {
					var resp WsResp
					resp.Id = req.Id
					resp.ErrMsg = err.Error()
					lock.Lock()
					defer lock.Unlock()
					return conn.WriteJSON(resp)
				})
				if err != nil {
					fmt.Printf("%+v %+v\n", req, err)
				}
				calculateBackupSizeCancelFunctions.Delete(req.Id)
			}()
		}
	}
}

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
	startTime   int64
	root        string
	lock        sync.Locker
	dbFileInfos []DbFileInfo
}

func (h *SizeWalkerHandlers) StackSizeHandler(size int) {
	h.StackSize = size
}

type DbFileInfo struct {
	time  int64
	path  string
	typ   int // 0: file 1: dir 2: root
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
			if filePath == h.root {
				typ = 2
			} else {
				typ = 1
			}
		}
		h.lock.Lock()
		defer h.lock.Unlock()
		h.dbFileInfos = append(h.dbFileInfos, DbFileInfo{
			time:  h.startTime,
			path:  filePath,
			typ:   typ,
			count: count,
			size:  size,
		})
	}
}

func (h *SizeWalkerHandlers) insertFiles(ctx context.Context) error {
	conn := h.db.getConn()
	defer h.db.putConn(conn)
	return dbBase.InsertBatch[DbFileInfo](ctx, conn, 32766, h.dbFileInfos, 7, getInsertItemQuery, func(args []interface{}, start int, item DbFileInfo) {
		args[start] = item.time
		args[start+1] = item.path
		args[start+2] = filepath.Dir(item.path)
		args[start+3] = filepath.Base(item.path)
		args[start+4] = item.typ
		args[start+5] = item.count
		args[start+6] = item.size
	})
}

func getInsertItemQuery(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
	INSERT INTO _file (
	    time,
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

func calculateBackupSize(ctx context.Context, db *DB, req WsReq, backupDir string, onResp func(finished bool, data interface{}) error, onErr func(error) error) error {
	if !filepath.IsAbs(backupDir) {
		return onErr(errors.New("请输入绝对路径"))
	}
	info, err := os.Lstat(backupDir)
	if err != nil {
		return onErr(err)
	}
	if !info.IsDir() {
		return onErr(errors.New("请输入一个目录"))
	}
	handlers := SizeWalkerHandlers{
		req:       req,
		tick:      time.Tick(time.Millisecond * 500),
		onResp:    onResp,
		db:        db,
		startTime: time.Now().UnixNano(),
		root:      backupDir,
		lock:      &sync.Mutex{},
	}
	err = onResp(false, handlers.FileSizeResp)
	if err != nil {
		return err
	}
	_, err = core.Walk[CountAndSize](ctx, backupDir, 15, &handlers)
	if err != nil {
		return onErr(err)
	}
	err = onResp(false, handlers.FileSizeResp)
	if err != nil {
		return err
	}
	err = handlers.insertFiles(ctx)
	if err != nil {
		return onErr(err)
	}
	return onResp(true, handlers.FileSizeResp)
}
