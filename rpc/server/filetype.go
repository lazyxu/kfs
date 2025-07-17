package server

import (
	"fmt"
	"io"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/dao"
)

func AnalyzeFileType(kfsCore *core.KFS, size int64, md5 string, sha256 string) (ft dao.FileType, err error) {
	header, err := getHeader(kfsCore, size, md5, sha256)
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
	fmt.Printf("%d %s %s %+v\n", size, md5, sha256, fileType)
	ft = NewFileType(fileType)
	return
}

func NewFileType(fileType types.Type) dao.FileType {
	return dao.FileType{
		Type:      fileType.MIME.Type,
		SubType:   fileType.MIME.Subtype,
		Extension: fileType.Extension,
	}
}

func getHeader(kfsCore *core.KFS, size int64, md5 string, sha256 string) (header []byte, err error) {
	rc, err := kfsCore.S2.ReadWithSize(size, md5, sha256)
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
