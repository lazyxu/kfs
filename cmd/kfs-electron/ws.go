package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lazyxu/kfs/cmd/kfs-electron/backup"
	"github.com/lazyxu/kfs/cmd/kfs-electron/db/gosqlite"
	"io"
	"log"
	"net/http"
	"sync"
)

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

func wsHandler(w http.ResponseWriter, r *http.Request, serverAddr string, db *gosqlite.DB) {
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

func (p *WsProcessor) process(ctx context.Context, db *gosqlite.DB) {
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
		case "upsertBackup":
			newCtx, cancelFunc := context.WithCancel(ctx)
			p.cancelFunctions.Store(req.Id, cancelFunc)
			data := req.Data.(map[string]interface{})
			name := data["name"].(string)
			description := data["description"].(string)
			srcPath := data["srcPath"].(string)
			driverName := data["driverName"].(string)
			dstPath := data["dstPath"].(string)
			concurrent := int(data["concurrent"].(float64))
			encoder := data["encoder"].(string)
			go func() {
				err := p.upsertBackup(newCtx, db, req, name, description, srcPath, driverName, dstPath, encoder, concurrent)
				if err != nil {
					fmt.Printf("%+v %+v\n", req, err)
				}
				p.cancelFunctions.Delete(req.Id)
			}()
		}
	}
}

func (p *WsProcessor) upsertBackup(ctx context.Context, db *gosqlite.DB, req WsReq, name, description, srcPath, driverName, dstPath, encoder string, concurrent int) error {
	err := backup.UpsertBackup(ctx, db, name, description, srcPath, driverName, dstPath, encoder, concurrent)
	if err != nil {
		return p.err(req, err)
	}
	return nil
}
