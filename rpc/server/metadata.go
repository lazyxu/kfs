package server

import (
	"context"
	"errors"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/dbBase"
)

func AnalyzeIfNoFileType(ctx context.Context, kfsCore *core.KFS, hash string) error {
	_, err := kfsCore.Db.GetFileType(ctx, hash)
	if errors.Is(err, dbBase.ErrNoRecords) {
		err = Analyze(ctx, kfsCore, hash)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func Analyze(ctx context.Context, kfsCore *core.KFS, hash string) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	ft, err := AnalyzeFileType(kfsCore, hash)
	if err != nil {
		return err
	}
	err = InsertExif(ctx, kfsCore, hash, &ft)
	if err != nil {
		return err
	}
	err = kfsCore.Db.UpsertFileType(ctx, hash, ft)
	if err != nil {
		return err
	}
	return nil
}
