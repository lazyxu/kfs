package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/silenceper/pool"

	"github.com/lazyxu/kfs/core"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Upload(ctx context.Context, branchName string, dstPath string, srcPath string, config core.UploadConfig) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	return withFS2[sqlite.Commit, sqlite.Branch](fs,
		func(c pb.KoalaFSClient) (commit sqlite.Commit, branch sqlite.Branch, err error) {
			srcPath, err = filepath.Abs(srcPath)
			if err != nil {
				return
			}
			idleTimeout := time.Second * 10
			p, err := pool.NewChannelPool(&pool.Config{
				InitialCap: 0,
				MaxCap:     config.Concurrent,
				MaxIdle:    config.Concurrent,
				Factory: func() (interface{}, error) {
					return net.Dial("tcp", "127.0.0.1:1124")
				},
				Close: func(i interface{}) error {
					return i.(net.Conn).Close()
				},
				Ping: func(i interface{}) error {
					conn := i.(net.Conn)
					_, err := conn.Write([]byte{0})
					if err != nil {
						return err
					}
					var pong uint8
					err = binary.Read(conn, binary.LittleEndian, &pong)
					if err != nil {
						return err
					}
					if pong != 0 {
						return fmt.Errorf("pong is %d, expected 0", pong)
					}
					return nil
				},
				IdleTimeout: idleTimeout,
			})
			if err != nil {
				return
			}
			defer p.Release()
			handlers := &uploadHandlers{
				c:             c,
				p:             p,
				uploadProcess: config.UploadProcess,
				encoder:       config.Encoder,
				ch:            make(chan *Process),
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				handlers.handleProcess(srcPath, config.Concurrent)
				wg.Done()
			}()
			fileResp, err := core.Walk[fileResp](ctx, srcPath, config.Concurrent, handlers)
			close(handlers.ch)
			wg.Wait()
			if err != nil {
				return
			}
			info, err := os.Stat(srcPath)
			if err != nil {
				return
			}
			fileOrDir := fileResp.fileOrDir
			modifyTime := uint64(info.ModTime().UnixNano())
			client, err := c.Upload(ctx)
			if err != nil {
				return
			}
			err = client.Send(&pb.UploadReq{
				Root: &pb.UploadReqRoot{
					BranchName: branchName,
					Path:       dstPath,
					DirItem: &pb.DirItem{
						Hash:       fileOrDir.Hash(),
						Name:       filepath.Base(dstPath),
						Mode:       uint64(info.Mode()),
						Size:       fileOrDir.Size(),
						Count:      fileOrDir.Count(),
						TotalCount: fileOrDir.TotalCount(),
						CreateTime: modifyTime,
						ModifyTime: modifyTime,
						ChangeTime: modifyTime,
						AccessTime: modifyTime,
					},
				},
			})
			if err == io.EOF {
				err = nil
			}
			if err != nil {
				return
			}
			resp, err := client.Recv()
			if err != nil {
				return
			}
			err = client.CloseSend()
			if err != nil {
				return
			}
			return sqlite.Commit{
					Id:   resp.Branch.CommitId,
					Hash: resp.Branch.Hash,
				}, sqlite.Branch{
					Name:     branchName,
					CommitId: resp.Branch.CommitId,
					Size:     resp.Branch.Size,
					Count:    resp.Branch.Count,
				}, nil
		})
}
