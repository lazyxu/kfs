package dao

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
)

type Dir struct {
	FileOrDir
	count      uint64
	totalCount uint64
}

func (dir Dir) Count() uint64 {
	return dir.count
}

func (dir Dir) TotalCount() uint64 {
	return dir.totalCount
}

func NewDir(hash string, size uint64, count uint64, totalCount uint64) Dir {
	return Dir{FileOrDir{hash, size}, count, totalCount}
}

func NewDirFromDirItem(item IDirItem) (Dir, error) {
	if !os.FileMode(item.GetMode()).IsDir() {
		return Dir{}, ErrExpectedDir
	}
	return Dir{FileOrDir{item.GetHash(), item.GetSize()}, item.GetCount(), item.GetTotalCount()}, nil
}

func writeMutil(w io.Writer, order binary.ByteOrder, data []any) {
	for _, v := range data {
		err := binary.Write(w, order, v)
		if err != nil {
			panic(err)
		}
	}
}

func (dir *Dir) Cal(dirItems []DirItem) {
	hash := sha256.New()
	err := binary.Write(hash, binary.LittleEndian, uint64(len(dirItems)))
	if err != nil {
		panic(err)
	}
	dir.size = 0
	dir.count = uint64(len(dirItems))
	for _, dirItem := range dirItems {
		writeMutil(hash, binary.LittleEndian, []any{
			[]byte(dirItem.Hash),
			[]byte(dirItem.Name),
			dirItem.Mode,
			dirItem.CreateTime,
			dirItem.ModifyTime,
			dirItem.ChangeTime,
			dirItem.AccessTime,
		})
		dir.size += dirItem.Size
		dir.totalCount += dirItem.TotalCount
	}
	dir.hash = hex.EncodeToString(hash.Sum(nil))
}
