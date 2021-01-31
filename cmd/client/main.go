package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/lazyxu/kfs/cmd/client/localfs"

	"github.com/lazyxu/kfs/cmd/client/rpc"
)

var homeDir string

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

func main() {
	rpc.RegisterInvokeMethod("ls", func(q *rpc.JsonQ) (interface{}, error) {
		dirname := q.RFindStringOrDefault("path", homeDir)
		infos, err := ioutil.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		return localfs.TransInfo(dirname, infos), nil
	})
	rpc.Register1ton("backup", func(ctx context.Context, q *rpc.JsonQ, ch chan<- interface{}) error {
		dirname := q.RFindStringOrDefault("params.path", homeDir)
		host := q.RFindStringOrDefault("params.host", "127.0.0.1:9092")
		b, err := localfs.NewBackUpCtx(ctx, host, dirname, []localfs.IgnoreRule{
			&localfs.FileNameIgnore{FileName: "Docker.raw"},
		})
		if err != nil {
			return err
		}
		ticker := time.NewTicker(500 * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					p := b.GetStatus()
					ch <- p
					if p.Done {
						ticker.Stop()
						return
					}
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()
		return b.Scan()
	})
	rpc.Start()
}
