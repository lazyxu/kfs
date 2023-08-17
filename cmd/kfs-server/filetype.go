package main

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/lazyxu/kfs/dao"
)

func AnalysisFileType(ctx context.Context) error {
	println("AnalysisFileType")
	// TODO: now remain is 0
	hashList, err := kfsCore.Db.ListExif(ctx)
	if err != nil {
		return err
	}
	for hash := range hashList {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
		}
		_, err = GetFileType(hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFileType(hash string) (e dao.Exif, err error) {
	header, err := getHeader(hash)
	if err != nil {
		return
	}
	typ, err := filetype.Get(header)
	if err != nil {
		return
	}
	fmt.Printf("%s %+v\n", hash, typ)
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
