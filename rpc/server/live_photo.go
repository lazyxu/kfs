package server

import (
	"context"
	"errors"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/dbBase"
	"path/filepath"
	"strings"
)

func UpsertLivePhoto(kfsCore *core.KFS, movHash string, driverId uint64, dirPath []string, movName string) error {
	ext := filepath.Ext(movName)
	// TODO: check .livp in baidu photo.
	if ext == ".MOV" {
		name := strings.TrimSuffix(movName, ext)
		heicPath := append(dirPath, name+".HEIC")
		heicFile, err := kfsCore.Db.GetDriverFile(context.TODO(), driverId, heicPath)
		if err != nil {
			if !errors.Is(err, dbBase.ErrNoSuchFileOrDir) {
				return err
			}
		} else {
		}
		jpgPath := append(dirPath, name+".JPG")
		jpgFile, err := kfsCore.Db.GetDriverFile(context.TODO(), driverId, jpgPath)
		if err != nil {
			if !errors.Is(err, dbBase.ErrNoSuchFileOrDir) {
				return err
			}
		}
		err = kfsCore.Db.UpsertLivePhoto(context.TODO(), movHash, heicFile.Hash, jpgFile.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}
