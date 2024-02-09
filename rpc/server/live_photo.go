package server

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/dbBase"
)

func UpsertLivePhoto(ctx context.Context, kfsCore *core.KFS, hash string, driverId uint64, dirPath []string, name string) error {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".mov" {
		prefix := strings.TrimSuffix(name, ext)
		heicPath := append(dirPath, prefix+".HEIC")
		heicFile, err1 := kfsCore.Db.GetDriverFile(ctx, driverId, heicPath)
		if err1 != nil {
			if !errors.Is(err1, dbBase.ErrNoSuchFileOrDir) {
				return err1
			}
		}
		jpgPath := append(dirPath, prefix+".JPG")
		jpgFile, err2 := kfsCore.Db.GetDriverFile(ctx, driverId, jpgPath)
		if err2 != nil {
			if !errors.Is(err1, dbBase.ErrNoSuchFileOrDir) {
				return err2
			}
		}

		if errors.Is(err1, dbBase.ErrNoSuchFileOrDir) && errors.Is(err2, dbBase.ErrNoSuchFileOrDir) {
			return nil
		}
		err := kfsCore.Db.UpsertLivePhoto(context.TODO(), hash, heicFile.Hash, jpgFile.Hash, "")
		if err != nil {
			return err
		}
	} else if ext == ".heic" {

	} else if ext == ".jpg" {

	} else if ext == ".livp" {
		_, _, err := UnzipIfLivp(ctx, kfsCore, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func getHash(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()
	hash := sha256.New()
	_, err = io.Copy(hash, rc)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func UnzipIfLivp(ctx context.Context, kfsCore *core.KFS, hash string) (movHash string, heicHash string, err error) {
	_, heicHash, err = kfsCore.Db.GetLivePhotoByLivp(ctx, hash)
	if err == nil {
		return
	}
	if !errors.Is(err, dbBase.ErrNoRecords) {
		return
	}
	return UnzipLivp(ctx, kfsCore, hash)
}

func UnzipLivp(ctx context.Context, kfsCore *core.KFS, hash string) (movHash string, heicHash string, err error) {
	select {
	case <-ctx.Done():
		err = context.Canceled
		return
	default:
	}

	src := kfsCore.S.GetFilePath(hash)

	r, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer r.Close()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) (string, error) {
		itemHash, err2 := getHash(f)
		if err2 != nil {
			return "", err2
		}
		rc, err2 := f.Open()
		if err2 != nil {
			return "", err2
		}
		defer rc.Close()
		_, err = kfsCore.S.Write(itemHash, func(w io.Writer, hasher io.Writer) (e error) {
			rr := io.TeeReader(rc, hasher)
			_, err = io.Copy(w, rr)
			return err
		})
		return itemHash, nil
	}
	for _, f := range r.File {
		var itemHash string
		itemHash, err = extractAndWriteFile(f)
		if err != nil {
			return
		}
		if strings.HasSuffix(f.Name, ".heic") {
			heicHash = itemHash
		} else if strings.HasSuffix(f.Name, ".mov") {
			movHash = itemHash
		} else {
			err = errors.New("invalid livp format, there are files other than .heic and .mov")
			return
		}
	}
	if heicHash == "" {
		err = errors.New("invalid livp format, missing .heic file")
		return
	}
	if movHash == "" {
		err = errors.New("invalid livp format, missing .mov file")
		return
	}
	err = kfsCore.Db.UpsertLivePhoto(context.TODO(), movHash, heicHash, "", hash)
	if err != nil {
		return
	}
	return
}
