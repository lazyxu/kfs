package backup

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/lazyxu/kfs/cmd/kfs-cli/upload"

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
		bar := upload.NewProcessBar(info, filename)
		defer bar.Close()
		relPath, err := filepath.Rel(v.backupPath, filename)
		if err != nil {
			return nil, err
		}
		hash, err := upload.SendHeader(bar, filename, info, relPath, func(metadata *pb.UploadReqMetadata) error {
			return v.client.Send(&pb.BackupReq{
				Header: &pb.BackupReqHeader{
					BranchName: v.branchName,
					Base:       "",
				},
				Metadata: metadata,
			})
		})
		if err != nil {
			return nil, err
		}
		err = upload.SendContent(bar, hash, filename, func(data []byte, isLast bool) error {
			return v.client.Send(&pb.BackupReq{Bytes: data, IsLast: isLast})
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
	}
	return nil, nil
}
