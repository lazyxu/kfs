package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/h2non/filetype/matchers"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

func PluginAnalyzeTranscodeNotAnalyzed(ctx context.Context, kfsCore *core.KFS, name string, file sqlc.File) error {
	if name != "" {
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".livp" && ext != ".heic" && ext != ".mov" {
			return nil
		}
	}
	err := AnalyzeTranscodeNotAnalyzed(ctx, kfsCore, file)
	if err != nil {
		return err
	}
	return nil
}

func AnalyzeTranscodeNotAnalyzed(ctx context.Context, kfsCore *core.KFS, file sqlc.File) error {
	// TODO: 过滤transcode目录下已有的文件
	// _, err = kfsCore.Q.GetLivePhotoByTranscode(ctx, sqlc.NInt64(file.ID))
	// if err == nil {
	// 	return
	// }
	// if !errors.Is(err, sql.ErrNoRows) {
	// 	return
	// }
	return AnalyzeTranscode(ctx, kfsCore, file)
}

func AnalyzeTranscode(ctx context.Context, kfsCore *core.KFS, file sqlc.File) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	id := file.ID
	size := file.Size
	md5 := file.Md5
	sha256 := file.Sha256
	m, err := kfsCore.Q.GetMetadata(ctx, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		AnalyzeMetadata(ctx, kfsCore, file)
		m, err = kfsCore.Q.GetMetadata(ctx, id)
	}
	if m.Extension == "zip" {
		transcodeFilePath := filepath.Join(kfsCore.TransCodeDir(), sha256+".jpg")
		_, err = os.Stat(transcodeFilePath)
		if os.IsNotExist(err) {
			_, heic, err := AnalyzeLivpNotAnalyzed(ctx, kfsCore, sqlc.File{
				ID:     id,
				Size:   size,
				Md5:    md5,
				Sha256: sha256,
			})
			if err != nil {
				return err
			}
			heicFile, err := kfsCore.Q.GetFileByID(ctx, heic)
			if err != nil {
				return err
			}
			rc, err := kfsCore.S2.ReadWithSize(heicFile.Size, heicFile.Md5, heicFile.Sha256)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			img = Orientation(img, m.Orientation.Int64)
			err = imaging.Save(img, transcodeFilePath)
			fmt.Printf("transcode livp to jpg: %s, err: %v\n", transcodeFilePath, err)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		return nil
	}
	if m.Extension == matchers.TypeHeif.Extension {
		transcodeFilePath := filepath.Join(kfsCore.TransCodeDir(), sha256+".jpg")
		_, err = os.Stat(transcodeFilePath)
		if os.IsNotExist(err) {
			rc, err := kfsCore.S2.ReadWithSize(size, md5, sha256)
			if err != nil {
				return err
			}
			defer rc.Close()
			img, err := goheif.Decode(rc) // CGO_ENABLED=1 https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release
			if err != nil {
				return err
			}
			img = Orientation(img, m.Orientation.Int64)
			err = imaging.Save(img, transcodeFilePath)
			fmt.Printf("transcode heic to jpg: %s, err: %v\n", transcodeFilePath, err)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		return nil
	} else if m.Extension == matchers.TypeMov.Extension {
		src := kfsCore.S2.GetFilePath(size, md5, sha256)
		transcodeFilePath := filepath.Join(kfsCore.TransCodeDir(), sha256+".mp4")
		_, err = os.Stat(transcodeFilePath)
		if os.IsNotExist(err) {
			err = SafeFFmpegRun(ctx, src, transcodeFilePath)
			fmt.Printf("transcode mov to mp4: %s, err: %v\n", transcodeFilePath, err)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		return nil
	}
	return err
}
