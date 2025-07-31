package core

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gofrs/flock"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/sqlc"
)

type KFS struct {
	Q            *sqlc.Queries
	Db           dao.Database
	S            dao.Storage
	storageDir   string
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

func createGlobalLockFile(root string) error {
	lockFile, err := os.Create(path.Join(root, lockFileName))
	if err != nil {
		return err
	}
	defer lockFile.Close()
	return nil
}

func (s *Storage2) GetFilePath(size int64, md5Hash string, sha256Hash string) string {
	modDir := strconv.Itoa(int(size % 1024)) // 新增取模目录层级
	itoaSize := strconv.Itoa(int(size))
	return filepath.Join(s.root, modDir, itoaSize+"_"+md5Hash+"_"+sha256Hash)
}

func ensureDirExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, dirPerm) // 递归创建目录[9,10](@ref)
	}
	return nil
}

func (s *Storage2) Create() error {
	_, err := os.Stat(s.root)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(path.Join(s.root), dirPerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = createGlobalLockFile(s.root)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (s *Storage2) AddReader(size int64, md5Hash string, sha256Hash string, getReaderIfNotExsit func() (io.Reader, error)) (bool, error) {
	lock := flock.New(path.Join(s.root, lockFileName))
	err := lock.Lock()
	if err != nil {
		return false, err
	}
	defer lock.Unlock()
	p := s.GetFilePath(size, md5Hash, sha256Hash)
	dirPath := filepath.Dir(p)
	if err := ensureDirExists(dirPath); err != nil {
		return false, err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0o200)
	if err != nil {
		if os.IsPermission(err) {
			// file exists
			return true, nil
		}
		return false, err
	}
	r, err := getReaderIfNotExsit()
	if err != nil {
		f.Close()
		os.Remove(p)
		return false, err
	}
	sha256Hasher := sha256.New()
	md5Hasher := md5.New()
	multiWriter := io.MultiWriter(sha256Hasher, md5Hasher, f)
	written, err := io.CopyN(multiWriter, r, size) // 一次写入，同时计算两个哈希，并且写入文件
	if err != nil {
		f.Close()
		os.Remove(p)
		return false, err
	}
	if written != size {
		f.Close()
		os.Remove(p)
		return false, fmt.Errorf("written size %d does not match file size %d", written, size)
	}
	sha256Hash2 := hex.EncodeToString(sha256Hasher.Sum(nil))
	if sha256Hash != sha256Hash2 {
		f.Close()
		os.Remove(p)
		return false, fmt.Errorf("invalid sha256: expected %s, actual %s", sha256Hash, sha256Hash2)
	}
	md5Hash2 := hex.EncodeToString(md5Hasher.Sum(nil))
	if md5Hash != md5Hash2 {
		f.Close()
		os.Remove(p)
		return false, fmt.Errorf("invalid md5: expected %s, actual %s", md5Hash, md5Hash2)
	}
	f.Close()
	err = os.Chmod(p, 0o400) // read only
	if err != nil {
		f.Close()
		os.Remove(p)
		return false, fmt.Errorf("failed to change file mode: %s", p)
	}
	return false, nil
}

func (s *Storage2) Add(md5Hash string, sha256Hash string, file *os.File) (bool, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}
	return s.AddReader(stat.Size(), md5Hash, sha256Hash, func() (io.Reader, error) {
		return file, nil
	})
}

type sizedReaderCloser struct {
	io.ReadSeekCloser
	size int64
}

func (rc sizedReaderCloser) Size() int64 {
	return rc.size
}

func (s *Storage2) ReadWithSize(size int64, md5Hash string, sha256Hash string) (dao.SizedReadCloser, error) {
	p := s.GetFilePath(size, md5Hash, sha256Hash)
	f, err := os.OpenFile(p, os.O_RDONLY, 0o200)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return sizedReaderCloser{f, info.Size()}, nil
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
	s2 := &Storage2{root: "tmp"}
	s2.Create()
	return &KFS{Db: db, S: s, S2: s2}, nil
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

func New3(q *sqlc.Queries, newDatabase func() (dao.Database, error), storageDir, thumbnailDir, transCodeDir string) (*KFS, error) {
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
	s2 := &Storage2{root: storageDir}
	s2.Create()
	return &KFS{Q: q, Db: db, S2: s2, storageDir: storageDir, thumbnailDir: thumbnailDir, transCodeDir: transCodeDir}, nil
}

func (fs *KFS) ThumbnailDir() string {
	return fs.thumbnailDir
}

func (fs *KFS) TransCodeDir() string {
	return fs.transCodeDir
}

func (fs *KFS) Close() error {
	var err1 error
	if fs.S != nil {
		err1 = fs.S.Close()
	}
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
