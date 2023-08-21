package main

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/lazyxu/kfs/dao"
)

func AnalysisFileType(ctx context.Context) error {
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
		fileType, err := GetFileType(hash)
		if err != nil {
			return err
		}
		_, err = kfsCore.Db.InsertFileType(ctx, hash, dao.FileType{
			Type:      fileType.MIME.Type,
			SubType:   fileType.MIME.Subtype,
			Extension: fileType.Extension,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFileType(hash string) (t types.Type, err error) {
	header, err := getHeader(hash)
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

func getHeader(hash string) (header []byte, err error) {
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
