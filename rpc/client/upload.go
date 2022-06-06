package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/silenceper/pool"

	"github.com/lazyxu/kfs/core"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/lazyxu/kfs/pb"
)

func (fs GRPCFS) Upload(ctx context.Context, branchName string, dstPath string, srcPath string, uploadProcess core.UploadProcess, concurrent int) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	return withFS2[sqlite.Commit, sqlite.Branch](fs,
		func(c pb.KoalaFSClient) (commit sqlite.Commit, branch sqlite.Branch, err error) {
			srcPath, err = filepath.Abs(srcPath)
			if err != nil {
				return
			}
			idleTimeout := time.Second * 10
			p, err := pool.NewChannelPool(&pool.Config{
				InitialCap: 0,
				MaxCap:     concurrent,
				MaxIdle:    concurrent,
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
			v := &uploadVisitor{
				c:             c,
				p:             p,
				uploadProcess: uploadProcess,
				concurrent:    concurrent,
			}
			v.connCh = make(chan net.Conn, concurrent)
			for i := 0; i < concurrent; i++ {
				var conn net.Conn
				conn, err = net.Dial("tcp", "127.0.0.1:1124")
				if err != nil {
					return
				}
				v.connCh <- conn
			}
			walker := storage.NewWalker[sqlite.FileOrDir](ctx, srcPath, v)
			scanResp, err := walker.Walk(concurrent > 1)
			if err != nil {
				return
			}
			info, err := os.Stat(srcPath)
			if err != nil {
				return
			}
			fileOrDir := scanResp.(sqlite.FileOrDir)
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
