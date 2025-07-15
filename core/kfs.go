package core

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/gofrs/flock"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/sqlc"
)

type KFS struct {
	Q            *sqlc.Queries
	Db           dao.Database
	S            dao.Storage
	thumbnailDir string
	transCodeDir string
	newStorage   func(root string) (dao.Storage, error)
	isSqlite     bool
	S2           *Storage2
}

type Storage2 struct {
	root string
}

const (
	dirPerm      = 0o700
	lockFileName = "index.lock"
)

func (s *Storage2) Add(md5Hash string, sha256Hash string, file *os.File) (bool, error) {
	lock := flock.New(path.Join(s.root, lockFileName))
	err := lock.Lock()
	if err != nil {
		return false, err
	}
	defer lock.Unlock()
	size, err := file.Stat()
	if err != nil {
		return false, err
	}
	itoaSize := strconv.Itoa(int(size.Size()))
	dirPath := path.Join(s.root, itoaSize)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, dirPerm)
		if err != nil {
			return false, err
		}
	} else if err != nil {
		return false, err
	}
	p := path.Join(dirPath, md5Hash+"_"+sha256Hash)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		if os.IsPermission(err) {
			// file exists
			return true, nil
		}
		return false, err
	}
	defer f.Close()
	sha256Hasher := sha256.New()
	md5Hasher := md5.New()
	multiWriter := io.MultiWriter(sha256Hasher, md5Hasher)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		os.Remove(p)
		return false, err
	}
	written, err := io.Copy(multiWriter, file) // 一次写入，同时计算两个哈希
	if err != nil {
		os.Remove(p)
		return false, err
	}
	if written != size.Size() {
		os.Remove(p)
		return false, fmt.Errorf("written size %d does not match file size %d", written, size.Size())
	}
	sha256Hash2 := hex.EncodeToString(sha256Hasher.Sum(nil))
	if sha256Hash != sha256Hash2 {
		os.Remove(p)
		return false, fmt.Errorf("invalid sha256: expected %s, actual %s", sha256Hash, sha256Hash2)
	}
	md5Hash2 := hex.EncodeToString(md5Hasher.Sum(nil))
	if md5Hash != md5Hash2 {
		os.Remove(p)
		return false, fmt.Errorf("invalid md5: expected %s, actual %s", md5Hash, md5Hash2)
	}
	err = os.Chmod(p, 0o400) // read only
	if err != nil {
		os.Remove(p)
		return false, fmt.Errorf("failed to change file mode: %s", p)
	}
	return false, nil
}

func New(newDatabase func() (dao.Database, error), newStorage func() (dao.Storage, error)) (*KFS, error) {
	s, err := newStorage()
	if err != nil {
		return nil, err
	}
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s, S2: &Storage2{root: "."}}, nil
}

func mkdirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(path, 0o700)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func New2(newDatabase func() (dao.Database, error), newStorage func() (dao.Storage, error), thumbnailDir, transCodeDir string) (*KFS, error) {
	s, err := newStorage()
	if err != nil {
		return nil, err
	}
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(thumbnailDir)
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(transCodeDir)
	if err != nil {
		return nil, err
	}
	return &KFS{Db: db, S: s, thumbnailDir: thumbnailDir, transCodeDir: transCodeDir}, nil
}

func New3(q *sqlc.Queries, newDatabase func() (dao.Database, error), newStorage func() (dao.Storage, error), thumbnailDir, transCodeDir string) (*KFS, error) {
	s, err := newStorage()
	if err != nil {
		return nil, err
	}
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(thumbnailDir)
	if err != nil {
		return nil, err
	}
	err = mkdirIfNotExist(transCodeDir)
	if err != nil {
		return nil, err
	}
	return &KFS{Q: q, Db: db, S: s, thumbnailDir: thumbnailDir, transCodeDir: transCodeDir}, nil
}

func (fs *KFS) ThumbnailDir() string {
	return fs.thumbnailDir
}

func (fs *KFS) TransCodeDir() string {
	return fs.transCodeDir
}

func (fs *KFS) Close() error {
	err1 := fs.S.Close()
	err2 := fs.Db.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (fs *KFS) Reset() error {
	err := fs.S.Remove()
	if err != nil {
		return err
	}
	err = fs.S.Create()
	if err != nil {
		return err
	}
	err = fs.Db.Remove()
	if err != nil {
		return err
	}
	err = fs.Db.Create()
	if err != nil {
		return err
	}
	return nil
}
