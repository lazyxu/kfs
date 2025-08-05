package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

func AnalyzeMetadataNotAnalyzed(ctx context.Context, kfsCore *core.KFS, file sqlc.File) error {
	_, err := kfsCore.Q.GetFileType(ctx, file.ID)
	if errors.Is(err, sql.ErrNoRows) {
		err = AnalyzeMetadata(ctx, kfsCore, file)
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

func AnalyzeMetadataByFileId(ctx context.Context, kfsCore *core.KFS, id int64) error {
	file, err := kfsCore.Q.GetFileByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("file with id %d not found", id)
		}
		return err
	}
	return AnalyzeMetadata(ctx, kfsCore, file)
}

func AnalyzeMetadata(ctx context.Context, kfsCore *core.KFS, file sqlc.File) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	ft, err := AnalyzeFileType(kfsCore, file.Size, file.Md5, file.Sha256)
	if err != nil {
		fmt.Printf("AnalyzeFileType error: %v\n", err)
		return err
	}
	err = InsertExif(ctx, kfsCore, file, &ft)
	if err != nil {
		fmt.Printf("InsertExif error: %v\n", err)
		return err
	}
	err = kfsCore.Q.UpsertFileType(ctx, sqlc.UpsertFileTypeParams{ID: file.ID, Type: ft.Type, SubType: ft.SubType, Extension: ft.Extension})
	if err != nil {
		return err
	}
	return nil
}
