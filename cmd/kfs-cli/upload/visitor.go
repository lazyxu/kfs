package upload

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/lazyxu/kfs/pb"
	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type uploadVisitor struct {
	storage.EmptyVisitor[sqlite.FileOrDir]
	client     pb.KoalaFS_BackupClient
	branchName string
	backupPath string
	bars       sync.Map
}

func (v *uploadVisitor) HasExit() bool {
	return true
}

func (v *uploadVisitor) Exit(ctx context.Context, filename string, info os.FileInfo, infos []os.FileInfo, rets []sqlite.FileOrDir) (sqlite.FileOrDir, error) {
	if info.Mode().IsRegular() {
		file, err := sqlite.NewFileByName(filename)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		bar := NewProcessBar(info, filename)
		defer bar.Close()
		hash, err := GetFileHash(bar, filename)
		if err != nil {
			return nil, err
		}
		err = SendContent(bar, hash, filename, func(data []byte, isFirst bool, isLast bool) error {
			if isFirst {
				return v.client.Send(&pb.BackupReq{
					File: &pb.BackupReqFile{
						Hash:        hash,
						Size:        uint64(info.Size()),
						Ext:         filepath.Ext(filename),
						Bytes:       data,
						IsLastChunk: isLast,
					},
				})
			}
			return v.client.Send(&pb.BackupReq{
				File: &pb.BackupReqFile{
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
		err := v.client.Send(&pb.BackupReq{
			Dir: &pb.BackupReqDir{DirItem: dirItems},
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
