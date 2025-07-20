package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

func AnalyzeIfNoFileType(ctx context.Context, kfsCore *core.KFS, id int64, size int64, md5 string, sha256 string) error {
	_, err := kfsCore.Q.GetFileType(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = Analyze(ctx, kfsCore, id, size, md5, sha256)
		if err != nil {
			return err
		}
	}
	if err != nil {
		fmt.Printf("GetFileType error: %v\n", err)
		return err
	}
	return nil
}

func Analyze(ctx context.Context, kfsCore *core.KFS, id int64, size int64, md5 string, sha256 string) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	ft, err := AnalyzeFileType(kfsCore, size, md5, sha256)
	if err != nil {
		fmt.Printf("AnalyzeFileType error: %v\n", err)
		return err
	}
	err = InsertExif(ctx, kfsCore, id, size, md5, sha256, &ft)
	if err != nil {
		fmt.Printf("InsertExif error: %v\n", err)
		return err
	}
	err = kfsCore.Q.UpsertFileType(ctx, sqlc.UpsertFileTypeParams{ID: id, Type: ft.Type, SubType: ft.SubType, Extension: ft.Extension})
	if err != nil {
		return err
	}
	return nil
}
