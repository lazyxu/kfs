package server

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/dbBase"
	"io"
	"os"
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
		err = kfsCore.Db.UpsertLivePhoto(context.TODO(), movHash, heicFile.Hash, jpgFile.Hash, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func Unzip(src, dest string) (files []string, err error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err = extractAndWriteFile(f)
		if err != nil {
			return
		}
		files = append(files, filepath.Join(dest, f.Name))
	}

	return
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

func UnzipLivp(ctx context.Context, kfsCore *core.KFS, hash string) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	src := kfsCore.S.GetFilePath(hash)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
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
	var heicHash string
	var movHash string
	for _, f := range r.File {
		itemHash, err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
		if strings.HasSuffix(f.Name, ".heic") {
			heicHash = itemHash
		} else if strings.HasSuffix(f.Name, ".mov") {
			movHash = itemHash
		} else {
			return errors.New("invalid livp format")
		}
	}
	if heicHash == "" || movHash == "" {
		return errors.New("invalid livp format")
	}
	err = kfsCore.Db.UpsertLivePhoto(context.TODO(), movHash, heicHash, "", hash)
	if err != nil {
		return err
	}
	return nil
}
