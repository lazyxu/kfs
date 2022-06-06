package client

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/silenceper/pool"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	c             pb.KoalaFSClient
	p             pool.Pool
	uploadProcess core.UploadProcess
	concurrent    int
	connCh        chan net.Conn
}

func (v *uploadVisitor) HasExit() bool {
	return true
}

func (v *uploadVisitor) Exit(ctx context.Context, filePath string, info os.FileInfo, infos []os.FileInfo, rets []sqlite.FileOrDir) (sqlite.FileOrDir, error) {
	if info.Mode().IsRegular() {
		v.uploadProcess = v.uploadProcess.New(int(info.Size()), filepath.Base(filePath))
		defer v.uploadProcess.Close()
		file, err := core.NewFileByName(v.uploadProcess, filePath)
		if err != nil {
			return nil, err
		}
		if v.concurrent > 1 {
			err = v.uploadFile(filePath, file.Hash(), file.Size())
			if err != nil {
				return nil, err
			}
			return file, nil
		}
		client, err := v.c.Upload(ctx)
		if err != nil {
			return nil, err
		}
		err = SendContent(v.uploadProcess, file.Hash(), filePath, func(data []byte, isFirst bool, isLast bool) error {
			if isFirst {
				return client.Send(&pb.UploadReq{
					File: &pb.UploadReqFile{
						Hash:        file.Hash(),
						Size:        uint64(info.Size()),
						Bytes:       data,
						IsLastChunk: isLast,
					},
				})
			}
			return client.Send(&pb.UploadReq{
				File: &pb.UploadReqFile{
					Bytes:       data,
					IsLastChunk: isLast,
				},
			})
		})
		if err != nil {
			return nil, err
		}
		_, err = client.Recv()
		if err != nil {
			return nil, err
		}
		return file, err
	} else if info.IsDir() {
		dirItems := make([]*pb.DirItem, len(infos))
		for i, info := range infos {
			if rets[i] == nil {
				continue
			}
			modifyTime := uint64(info.ModTime().UnixNano())
			dirItems[i] = &pb.DirItem{
				Hash:       rets[i].Hash(),
				Name:       info.Name(),
				Mode:       uint64(info.Mode()),
				Size:       rets[i].Size(),
				Count:      rets[i].Count(),
				TotalCount: rets[i].TotalCount(),
				CreateTime: modifyTime,
				ModifyTime: modifyTime,
				ChangeTime: modifyTime,
				AccessTime: modifyTime,
			}
		}
		client, err := v.c.Upload(ctx)
		if err != nil {
			return nil, err
		}
		err = client.Send(&pb.UploadReq{
			Dir: &pb.UploadReqDir{DirItem: dirItems},
		})
		resp, err := client.Recv()
		if err != nil {
			return nil, err
		}
		err = client.CloseSend()
		if err != nil {
			return nil, err
		}
		dir := sqlite.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return dir, nil
	}
	return nil, nil
}
