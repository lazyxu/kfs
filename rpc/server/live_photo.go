package server

import (
	"archive/zip"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/db/sqlc"
)

func PluginUnzipIfLivp(ctx context.Context, kfsCore *core.KFS, hash string, name string) error {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".livp" {
		return nil
	}
	_, _, err := UnzipIfLivp(ctx, kfsCore, sqlc.File{})
	if err != nil {
		return err
	}
	return nil
}

func getHash(f *zip.File) (string, string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", "", err
	}
	defer rc.Close()

	// 创建两个哈希计算器
	hashSHA256 := sha256.New()
	hashMD5 := md5.New()

	// 创建管道：将数据同时分流到 SHA256 和 MD5
	pipeReader, pipeWriter := io.Pipe()
	tee := io.TeeReader(rc, pipeWriter) // 分流到管道写入端

	// 启动协程：从管道读取端计算 MD5
	var wg sync.WaitGroup
	wg.Add(1)
	var md5Sum string
	var md5Err error
	go func() {
		defer wg.Done()
		_, md5Err = io.Copy(hashMD5, pipeReader)
		md5Sum = hex.EncodeToString(hashMD5.Sum(nil))
	}()

	// 主线程：从 tee 读取并计算 SHA256
	_, err = io.Copy(hashSHA256, tee) // 读取时同步写入管道
	pipeWriter.Close()                // 确保管道写入结束
	if err != nil {
		return "", "", err
	}

	// 等待 MD5 计算完成
	wg.Wait()
	if md5Err != nil {
		return "", "", md5Err
	}

	shaSum := hex.EncodeToString(hashSHA256.Sum(nil))
	return md5Sum, shaSum, nil
}

func UnzipIfLivp(ctx context.Context, kfsCore *core.KFS, file sqlc.File) (mov int64, heic int64, err error) {
	_, err = kfsCore.Q.GetLivePhotoByLivp(ctx, sqlc.NInt64(file.ID))
	if err == nil {
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return
	}
	return UnzipLivp(ctx, kfsCore, file)
}

func UnzipLivp(ctx context.Context, kfsCore *core.KFS, file sqlc.File) (int64, int64, error) {
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
	otherCnt := 0
	for _, zf := range r.File {
		if strings.HasSuffix(zf.Name, ".heic") {
			heicCnt++
		} else if strings.HasSuffix(zf.Name, ".mov") {
			movCnt++
		} else {
			otherCnt++
		}
	}
	if otherCnt > 0 || heicCnt != 1 || movCnt != 1 {
		err = errors.New("invalid livp format, there are files other than .heic and .mov")
		return 0, 0, err
	}
	for _, zf := range r.File {
		fmt.Printf("\tzf: %+v\n", zf)
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
