package server

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

func PluginAnalyzeLivpNotAnalyzed(ctx context.Context, kfsCore *core.KFS, hash string, name string) error {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".livp" {
		return nil
	}
	_, _, err := AnalyzeLivpNotAnalyzed(ctx, kfsCore, sqlc.File{})
	if err != nil {
		return err
	}
	return nil
}

func AnalyzeLivpNotAnalyzed(ctx context.Context, kfsCore *core.KFS, file sqlc.File) (mov int64, heic int64, err error) {
	_, err = kfsCore.Q.GetLivePhotoByLivp(ctx, sqlc.NInt64(file.ID))
	if err == nil {
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return
	}
	return AnalyzeLivp(ctx, kfsCore, file)
}

func AnalyzeLivp(ctx context.Context, kfsCore *core.KFS, file sqlc.File) (int64, int64, error) {
	select {
	case <-ctx.Done():
		return 0, 0, context.Canceled
	default:
	}

	src := kfsCore.S2.GetFilePath(file.Size, file.Md5, file.Sha256)

	r, err := zip.OpenReader(src)
	if err != nil {
		return 0, 0, err
	}
	defer r.Close()

	var heic, mov int64
	fmt.Printf("UnzipLivp livp: %+v\n", file)
	heicCnt := 0
	movCnt := 0
	for _, zf := range r.File {
		fmt.Printf("\tzf: %+v\n", zf)
		if strings.HasSuffix(zf.Name, ".heic") {
			heicCnt++
		} else if strings.HasSuffix(zf.Name, ".mov") {
			movCnt++
		} else {
			return 0, 0, errors.New("invalid livp format, there are files other than .heic and .mov")
		}
	}
	if heicCnt != 1 || movCnt != 1 {
		err = errors.New("invalid livp format, there are files other than .heic and .mov")
		return 0, 0, err
	}
	for _, zf := range r.File {
		md5, sha256, err := getHash(zf)
		if err != nil {
			return 0, 0, err
		}
		rc, err := zf.Open()
		if err != nil {
			return 0, 0, err
		}
		defer rc.Close()
		size := zf.FileInfo().Size()
		_, err = kfsCore.S2.AddReader(size, md5, sha256, rc)
		if err != nil {
			return 0, 0, err
		}
		f := sqlc.UpsertFileParams{
			Size:   size,
			Md5:    md5,
			Sha256: sha256,
		}
		res, err := kfsCore.Q.UpsertFile(ctx, f)
		if err != nil {
			return 0, 0, err
		}
		fileID, err := res.LastInsertId()
		if err != nil {
			return 0, 0, err
		}
		if strings.HasSuffix(zf.Name, ".heic") {
			heic = fileID
		} else if strings.HasSuffix(zf.Name, ".mov") {
			mov = fileID
		} else {
			panic("invalid livp format, there are files other than .heic and .mov")
		}
	}
	if heic == 0 {
		err = errors.New("invalid livp format, missing .heic file")
		return 0, 0, err
	}
	if mov == 0 {
		err = errors.New("invalid livp format, missing .mov file")
		return 0, 0, err
	}
	err = kfsCore.Q.UpsertLivePhoto(context.TODO(), sqlc.UpsertLivePhotoParams{
		Mov:  mov,
		Heic: sqlc.NInt64(heic),
		Jpg:  sqlc.NInt64Null(),
		Livp: sqlc.NInt64(file.ID),
	})
	if err != nil {
		return 0, 0, err
	}
	return mov, heic, err
}
