package core

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

func GetHash(r io.Reader) (string, string, error) {
	// 创建两个哈希计算器
	hashSHA256 := sha256.New()
	hashMD5 := md5.New()

	// 创建管道：将数据同时分流到 SHA256 和 MD5
	pipeReader, pipeWriter := io.Pipe()
	tee := io.TeeReader(r, pipeWriter) // 分流到管道写入端

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
	_, err := io.Copy(hashSHA256, tee) // 读取时同步写入管道
	pipeWriter.Close()                 // 确保管道写入结束
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

func GetHashFromZipFile(f *zip.File) (string, string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", "", err
	}
	defer rc.Close()
	return GetHash(rc)
}

func GetHashFromFile(f *os.File) (string, string, error) {
	origOffset, seekErr := f.Seek(0, io.SeekCurrent)
	if seekErr != nil {
		return "", "", seekErr
	}

	if _, seekErr = f.Seek(0, io.SeekStart); seekErr != nil {
		return "", "", seekErr
	}

	defer func() {
		// 恢复原始文件指针位置
		f.Seek(origOffset, io.SeekStart)
	}()
	return GetHash(f)
}
