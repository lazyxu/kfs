package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/lazyxu/kfs/pb"
)

func SendHeader(bar io.Writer, filename string, info os.FileInfo, p string, fn func(*pb.UploadReqMetadata) error) (string, error) {
	hash, err := getFileHash(bar, filename)
	if err != nil {
		return "", err
	}
	modifyTime := uint64(info.ModTime().UnixNano())
	err = fn(&pb.UploadReqMetadata{
		Path:       p,
		Hash:       hash,
		Mode:       uint64(info.Mode()),
		Size:       uint64(info.Size()),
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	})
	if err != nil {
		return "", err
	}
	return hash, err
}

func getFileHash(bar io.Writer, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := sha256.New()
	w := io.MultiWriter(hash, bar)
	_, err = io.Copy(w, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
