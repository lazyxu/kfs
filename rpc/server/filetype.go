package server

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
	"io"
)

func AnalyzeFileType(kfsCore *core.KFS, hash string) (ft dao.FileType, err error) {
	header, err := getHeader(kfsCore, hash)
	if err == io.EOF {
		ft = NewFileType(filetype.Unknown)
		err = nil
		return
	}
	if err != nil {
		return
	}
	fileType, err := filetype.Get(header)
	if err != nil {
		return
	}
	fmt.Printf("%s %+v\n", hash, fileType)
	ft = NewFileType(fileType)
	return
}

func InsertFileType(ctx context.Context, kfsCore *core.KFS, hash string, ft dao.FileType) (err error) {
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
