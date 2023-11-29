package dao

import (
	"context"
	"os"
)

type Database interface {
	IsSqlite() bool
	DataSourceName() string
	Size() (int64, error)
	Remove() error
	Create() error
	Close() error

	ResetBranch(ctx context.Context, branchName string) error
	NewBranch(ctx context.Context, branchName string) (exist bool, err error)
	DeleteBranch(ctx context.Context, branchName string) error
	BranchInfo(ctx context.Context, branchName string) (branch Branch, err error)
	BranchList(ctx context.Context) (branches []IBranch, err error)

	WriteCommit(ctx context.Context, commit *Commit) error

	WriteDir(ctx context.Context, dirItems []DirItem) (dir Dir, err error)
	RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit Commit, branch Branch, err error)

	WriteFile(ctx context.Context, file File) error
	UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item DirItem) (commit Commit, branch Branch, err error)
	UpsertDirItems(ctx context.Context, branchName string, splitPath []string, items []DirItem) (commit Commit, branch Branch, err error)
	GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error)

	List(ctx context.Context, branchName string, splitPath []string) (dirItems []DirItem, err error)
	ListByHash(ctx context.Context, hash string) (dirItems []DirItem, err error)

	Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []DirItem, err error)
	Open2(ctx context.Context, branchName string, splitPath []string) (dirItem DirItem, dirItems []DirItem, err error)

	FileCount(ctx context.Context) (int, error)
	DirCount(ctx context.Context) (int, error)
	DirItemCount(ctx context.Context) (int, error)
	BranchCount(ctx context.Context) (int, error)

	InsertDevice(ctx context.Context, name string, os string) (int64, error)
	DeleteDevice(ctx context.Context, deviceId uint64) error
	ListDevice(ctx context.Context) (devices []Device, err error)

	InsertDriver(ctx context.Context, driverName string, description string) (exist bool, err error)
	InsertDriverBaiduPhoto(ctx context.Context, driverName string, description string, accessToken string, refreshToken string) (exist bool, err error)
	InsertDriverLocalFile(ctx context.Context, driverName string, description string, deviceId uint64, srcPath string, ignores string, encoder string) (exist bool, err error)
	UpdateDriverSync(ctx context.Context, driverId uint64, sync bool, h int64, m int64) error
	UpdateDriverLocalFile(ctx context.Context, driverId uint64, srcPath, ignores, encoder string) error
	ResetDriver(ctx context.Context, driverId uint64) error
	DeleteDriver(ctx context.Context, driverId uint64) error
	ListDriver(ctx context.Context) (drivers []Driver, err error)
	GetDriver(ctx context.Context, driverId uint64) (driver Driver, err error)
	GetDriverToken(ctx context.Context, driverId uint64) (driver Driver, err error)
	GetDriverSync(ctx context.Context, driverId uint64) (driver Driver, err error)
	ListCloudDriverSync(ctx context.Context) (drivers []Driver, err error)
	ListLocalFileDriver(ctx context.Context, deviceId uint64) (drivers []Driver, err error)
	GetDriverLocalFile(ctx context.Context, driverId uint64) (driver Driver, err error)

	GetDriverFileSize(ctx context.Context, driverId uint64) (n uint64, err error)
	GetDriverFileCount(ctx context.Context, driverId uint64) (n uint64, err error)
	GetDriverDirCount(ctx context.Context, driverId uint64) (n uint64, err error)

	InsertFile(ctx context.Context, hash string, size uint64) error
	InsertFileMd5(ctx context.Context, hash string, hashMd5 string) error
	ListFileMd5(ctx context.Context, md5List []string) (m map[string]string, err error)
	SumFileSize(ctx context.Context) (size uint64, err error)

	UpsertDriverFile(ctx context.Context, f DriverFile) error
	ListDriverFile(ctx context.Context, driverId uint64, filePath []string) (files []DriverFile, err error)
	GetDriverFile(ctx context.Context, driverId uint64, splitPath []string) (file DriverFile, err error)
	ListDriverFileByHash(ctx context.Context, hash string) (files []DriverFile, err error)
	CheckExists(ctx context.Context, driverId uint64, dirPath []string, checks []DirItemCheck, hashList []string) error

	InsertHeightWidth(ctx context.Context, hash string, hw HeightWidth) (exist bool, err error)
	InsertNullVideoMetadata(ctx context.Context, hash string) (exist bool, err error)
	InsertVideoMetadata(ctx context.Context, hash string, m VideoMetadata) (exist bool, err error)

	InsertNullExif(ctx context.Context, hash string) (exist bool, err error)
	InsertExif(ctx context.Context, hash string, e Exif) (exist bool, err error)
	ListExpectExif(ctx context.Context) (hashList []string, err error)
	ListExpectExifCb(ctx context.Context, cb func(hash string)) (err error)
	ListExif(ctx context.Context) (exifMap map[string]Exif, err error)
	ListMetadata(ctx context.Context) (list []Metadata, err error)
	GetMetadata(ctx context.Context, hash string) (Metadata, error)

	ListFile(ctx context.Context) (hashList []string, err error)

	InsertFileType(ctx context.Context, hash string, t FileType) (exist bool, err error)
	ListExpectFileType(ctx context.Context) (hashList []string, err error)
	GetFileType(ctx context.Context, hash string) (fileType FileType, err error)

	UpsertLivePhoto(ctx context.Context, movHash string, heicHash string, jpgHash string) error
}

func DatabaseNewFunc(dataSourceName string, newDB func(dataSourceName string) (Database, error)) func() (Database, error) {
	return func() (Database, error) {
		return newDB(dataSourceName)
	}
}
