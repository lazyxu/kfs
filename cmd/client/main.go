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
		b := localfs.NewBackUpCtx(ctx, dirname, []localfs.IgnoreRule{
			&localfs.FileNameIgnore{FileName: "Docker.raw"},
		})
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
		b.Scan()
		return nil
	})
	rpc.Start()
}
