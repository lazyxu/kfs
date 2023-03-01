package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lazyxu/kfs/core"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	cmd.PersistentFlags().String(PortStr, "0", "local web server port")
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
	//defer func() {
	//	if err := recover(); err != nil {
	//		println("recover", err)
	//		conn.Close()
	//	}
	//}()

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
		case "calculateBackupSize":
			go func() {
				err := calculateBackupSize(ctx, req.Data.(map[string]interface{})["backupDir"].(string), func(finished bool, data interface{}) error {
					var resp WsResp
					resp.Id = req.Id
					resp.Finished = finished
					resp.Data = data
					return conn.WriteJSON(resp)
				}, func(err error) error {
					var resp WsResp
					resp.Id = req.Id
					resp.ErrMsg = err.Error()
					return conn.WriteJSON(resp)
				})
				if err != nil {
					fmt.Printf("%+v %+v\n", req, err)
				}
			}()
		}
	}
}

type SizeWalkerHandlers struct {
	core.DefaultWalkHandlers[int64]
	onResp func(finished bool, data interface{}) error
	tick   <-chan time.Time
	total  int64
}

func (h SizeWalkerHandlers) FileHandler(ctx context.Context, index int, filePath string, info os.FileInfo, children []int64) int64 {
	var size int64
	if !info.IsDir() {
		h.total += info.Size()
		size = info.Size()
		return size
	}
	for _, child := range children {
		size += child
	}

	select {
	case <-h.tick:
		err := h.onResp(false, h.total)
		if err != nil {
			panic(err)
		}
	case <-ctx.Done():
		return size
	}
	return size
}

func calculateBackupSize(ctx context.Context, backupDir string, onResp func(finished bool, data interface{}) error, onErr func(error) error) error {
	if !filepath.IsAbs(backupDir) {
		return onErr(errors.New("请输入绝对路径"))
	}
	err := onResp(false, 0)
	if err != nil {
		return err
	}
	handlers := SizeWalkerHandlers{
		tick:   time.Tick(time.Millisecond * 500),
		onResp: onResp,
	}
	resp, err := core.Walk[int64](ctx, backupDir, 15, handlers)
	if err != nil {
		return onErr(errors.New("请输入绝对路径"))
	}
	return onResp(true, resp)
}
