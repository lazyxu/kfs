package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
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

	db, err := NewDb("electron.db")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, serverAddr, db)
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

type WsProcessor struct {
	conn            *websocket.Conn
	cancelFunctions sync.Map
	lock            sync.Mutex
}

func wsHandler(w http.ResponseWriter, r *http.Request, serverAddr string, db *DB) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	p := WsProcessor{
		conn: c,
	}
	p.process(r.Context(), db)
}

type WsReq struct {
	Type string      `json:"type"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}

type WsResp struct {
	Id       string      `json:"id"`
	Finished bool        `json:"finished"`
	Data     interface{} `json:"data"`
	ErrMsg   string      `json:"errMsg,omitempty"`
}

func (p *WsProcessor) ok(req WsReq, finished bool, data interface{}) error {
	var resp WsResp
	resp.Id = req.Id
	resp.Finished = finished
	resp.Data = data
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.conn.WriteJSON(resp)
}

func (p *WsProcessor) err(req WsReq, err error) error {
	var resp WsResp
	resp.Id = req.Id
	resp.Finished = true
	resp.ErrMsg = err.Error()
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.conn.WriteJSON(resp)
}

func (p *WsProcessor) process(ctx context.Context, db *DB) {
	println(p.conn.RemoteAddr().String(), "Process")
	defer func() {
		p.cancelFunctions.Range(func(key, value any) bool {
			cancelFunc, ok := p.cancelFunctions.Load(key)
			if !ok {
				return true
			}
			cancelFunc.(context.CancelFunc)()
			return true
		})
	}()
	//defer func() {
	//	if err := recover(); err != nil {
	//		println("recover", err)
	//		conn.Close()
	//	}
	//}()

	for {
		print(p.conn.RemoteAddr().String(), " ReadJSON ")
		var req WsReq
		err := p.conn.ReadJSON(&req)
		if err == io.EOF || websocket.IsUnexpectedCloseError(err) {
			p.conn.Close()
			println()
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", req)
		switch req.Type {
		case "scan.cancel":
			fallthrough
		case "fastScan.cancel":
			fallthrough
		case "cancel":
			cancelFunc, ok := p.cancelFunctions.Load(req.Id)
			if !ok {
				p.ok(req, true, nil)
				continue
			}
			cancelFunc.(context.CancelFunc)()
			p.cancelFunctions.Delete(req.Id)
		case "scan":
			newCtx, cancelFunc := context.WithCancel(ctx)
			p.cancelFunctions.Store(req.Id, cancelFunc)
			data := req.Data.(map[string]interface{})
			srcPath := data["srcPath"].(string)
			record := data["record"].(bool)
			concurrent := int(data["concurrent"].(float64))
			go func() {
				var err error
				if !record {
					err = p.fastScan(newCtx, req, srcPath, concurrent)
				} else {
					err = p.scan(newCtx, db, req, srcPath, concurrent)
				}
				if err != nil {
					fmt.Printf("%+v %+v\n", req, err)
				}
				p.cancelFunctions.Delete(req.Id)
			}()
		case "fastBackup":
			newCtx, cancelFunc := context.WithCancel(ctx)
			p.cancelFunctions.Store(req.Id, cancelFunc)
			data := req.Data.(map[string]interface{})
			srcPath := data["srcPath"].(string)
			serverAddr := data["serverAddr"].(string)
			branchName := data["branchName"].(string)
			dstPath := data["dstPath"].(string)
			concurrent := int(data["concurrent"].(float64))
			encoder := data["encoder"].(string)
			go func() {
				err := p.fastBackup(newCtx, req, srcPath, serverAddr, branchName, dstPath, concurrent, encoder)
				if err != nil {
					fmt.Printf("%+v %+v\n", req, err)
				}
				p.cancelFunctions.Delete(req.Id)
			}()
		}
	}
}
