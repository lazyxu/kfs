package noncgo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"io"
)

type Dir struct {
	fileOrDir
	count uint64
}

func (i Dir) Count() uint64 {
	return i.count
}

// https://zhuanlan.zhihu.com/p/343682839
type DirItem struct {
	Hash        string
	Name        string
	Mode        uint64
	Size        uint64
	Count       uint64
	CreateTime  uint64 // linux does not support it.
	ModifyTime  uint64
	ChangeTime  uint64 // windows does not support it.
	AccessTime  uint64
	OldItemHash string
}

func NewDirItem(fileOrDir FileOrDir, name string, mode uint64, createTime uint64, modifyTime uint64, changeTime uint64, accessTime uint64, oldItemHash string) DirItem {
	return DirItem{fileOrDir.Hash(), name, mode, fileOrDir.Size(), fileOrDir.Count(), createTime, modifyTime, changeTime, accessTime, oldItemHash}
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
	dir.count = 0
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
		dir.count += dirItem.Count
	}
	dir.hash = hex.EncodeToString(hash.Sum(nil))
}

func (db *SqliteNonCgoDB) WriteDir(ctx context.Context, dirItems []DirItem) (dir Dir, err error) {
	dir.Cal(dirItems)
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			println(err.Error())
			err = tx.Rollback()
			if err != sql.ErrTxDone {
				return
			}
		}
	}()
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
	}()
	_, err = tx.ExecContext(ctx, `
	INSERT INTO dir VALUES (?, ?, ?);
	`, dir.hash, dir.size, dir.count)
	if err != nil {
		return
	}
	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO dirItem (
		hash,
		itemHash,
		itemName,
		itemMode,
		itemSize,
		itemCount,
		itemCreateTime,
		itemModifyTime,
		itemChangeTime,
		itemAccessTime,
		oldItemHash
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = stmt.Close()
		}
	}()
	for _, dirItem := range dirItems {
		_, err = stmt.ExecContext(ctx,
			dir.hash,
			dirItem.Hash,
			dirItem.Name,
			dirItem.Mode,
			dirItem.Size,
			dirItem.Count,
			dirItem.CreateTime,
			dirItem.ModifyTime,
			dirItem.ChangeTime,
			dirItem.AccessTime,
			dirItem.OldItemHash)
		if err != nil {
			return
		}
	}
	return
}
