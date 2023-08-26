package server

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
)

func AnalysisFileType(ctx context.Context, kfsCore *core.KFS) error {
	println("AnalysisFileType")
	// TODO: now remain is 0
	hashList, err := kfsCore.Db.ListFile(ctx)
	if err != nil {
		return err
	}
	for _, hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		_, err = InsertFileType(ctx, kfsCore, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertFileType(ctx context.Context, kfsCore *core.KFS, hash string) (ft dao.FileType, err error) {
	header, err := getHeader(kfsCore, hash)
	if err != nil {
		return
	}
	fileType, err := filetype.Get(header)
	if err != nil {
		return
	}
	fmt.Printf("%s %+v\n", hash, fileType)
	ft = NewFileType(fileType)
	_, err = kfsCore.Db.InsertFileType(ctx, hash, ft)
	if err != nil {
		return
	}
	return
}

func NewFileType(fileType types.Type) dao.FileType {
	return dao.FileType{
		Type:      fileType.MIME.Type,
		SubType:   fileType.MIME.Subtype,
		Extension: fileType.Extension,
	}
}
func GetFileType(kfsCore *core.KFS, hash string) (t types.Type, err error) {
	header, err := getHeader(kfsCore, hash)
	if err != nil {
		return
	}
	t, err = filetype.Get(header)
	if err != nil {
		return
	}
	fmt.Printf("%s %+v\n", hash, t)
	return
}

func getHeader(kfsCore *core.KFS, hash string) (header []byte, err error) {
	rc, err := kfsCore.S.ReadWithSize(hash)
	if err != nil {
		return
	}
	defer rc.Close()
	header = make([]byte, 261)
	_, err = rc.Read(header)
	if err != nil {
		return
	}
	return
}
