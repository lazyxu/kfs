package client

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	client        pb.KoalaFS_UploadClient
	uploadProcess core.UploadProcess
	concurrent    bool
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
		if v.concurrent {
			err = uploadFile(filePath, file.Hash(), file.Size())
			if err != nil {
				return nil, err
			}
			return file, nil
		}
		err = SendContent(v.uploadProcess, file.Hash(), filePath, func(data []byte, isFirst bool, isLast bool) error {
			if isFirst {
				return v.client.Send(&pb.UploadReq{
					File: &pb.UploadReqFile{
						Hash:        file.Hash(),
						Size:        uint64(info.Size()),
						Bytes:       data,
						IsLastChunk: isLast,
					},
				})
			}
			return v.client.Send(&pb.UploadReq{
				File: &pb.UploadReqFile{
					Bytes:       data,
					IsLastChunk: isLast,
				},
			})
		})
		if err != nil {
			return nil, err
		}
		_, err = v.client.Recv()
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
		err := v.client.Send(&pb.UploadReq{
			Dir: &pb.UploadReqDir{DirItem: dirItems},
		})
		resp, err := v.client.Recv()
		if err != nil {
			return nil, err
		}
		dir := sqlite.NewDir(resp.Dir.Hash, resp.Dir.Size, resp.Dir.Count, resp.Dir.TotalCount)
		return dir, nil
	}
	return nil, nil
}
