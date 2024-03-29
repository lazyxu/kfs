package dao

type DirItemCheck struct {
	Name       string
	Size       uint64
	ModifyTime uint64
}

type IDriverFile interface {
	GetDriverId() uint64
	GetDirPath() []string
	GetName() string
	GetVersion() uint64
	GetHash() string
	GetMode() uint64
	GetSize() uint64
	GetCreateTime() uint64
	GetModifyTime() uint64
	GetChangeTime() uint64
	GetAccessTime() uint64
}

// https://zhuanlan.zhihu.com/p/343682839
type DriverFile struct {
	DriverId       uint64   `json:"driverId"`
	DriverName     string   `json:"driverName"`
	DirPath        []string `json:"dirPath"`
	Name           string   `json:"name"`
	Version        uint64   `json:"version"` // TODO: REMOVE IT!
	Hash           string   `json:"hash"`
	Mode           uint64   `json:"mode"`
	Size           uint64   `json:"size"`
	CreateTime     uint64   `json:"createTime"` // linux does not support it.
	ModifyTime     uint64   `json:"modifyTime"`
	ChangeTime     uint64   `json:"changeTime"` // windows does not support it.
	AccessTime     uint64   `json:"accessTime"`
	UploadDeviceId string   `json:"uploadDeviceId"`
	UploadTime     uint64   `json:"uploadTime"`
}

func (d DriverFile) GetDriverId() uint64 {
	return d.DriverId
}

func (d DriverFile) GetDirPath() []string {
	return d.DirPath
}

func (d DriverFile) GetName() string {
	return d.Name
}

func (d DriverFile) GetHash() string {
	return d.Hash
}

func (d DriverFile) GetMode() uint64 {
	return d.Mode
}

func (d DriverFile) GetSize() uint64 {
	return d.Size
}

func (d DriverFile) GetCreateTime() uint64 {
	return d.CreateTime
}

func (d DriverFile) GetModifyTime() uint64 {
	return d.ModifyTime
}

func (d DriverFile) GetChangeTime() uint64 {
	return d.ChangeTime
}

func (d DriverFile) GetAccessTime() uint64 {
	return d.AccessTime
}
