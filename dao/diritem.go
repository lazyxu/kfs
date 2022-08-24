package dao

type IDirItem interface {
	GetHash() string
	GetName() string
	GetMode() uint64
	GetSize() uint64
	GetCount() uint64
	GetTotalCount() uint64
	GetCreateTime() uint64
	GetModifyTime() uint64
	GetChangeTime() uint64
	GetAccessTime() uint64
}

// https://zhuanlan.zhihu.com/p/343682839
type DirItem struct {
	Hash       string
	Name       string
	Mode       uint64
	Size       uint64
	Count      uint64
	TotalCount uint64
	CreateTime uint64 // linux does not support it.
	ModifyTime uint64
	ChangeTime uint64 // windows does not support it.
	AccessTime uint64
}

func (d DirItem) GetHash() string {
	return d.Hash
}

func (d DirItem) GetName() string {
	return d.Name
}

func (d DirItem) GetMode() uint64 {
	return d.Mode
}

func (d DirItem) GetSize() uint64 {
	return d.Size
}

func (d DirItem) GetCount() uint64 {
	return d.Count
}

func (d DirItem) GetTotalCount() uint64 {
	return d.TotalCount
}

func (d DirItem) GetCreateTime() uint64 {
	return d.CreateTime
}

func (d DirItem) GetModifyTime() uint64 {
	return d.ModifyTime
}

func (d DirItem) GetChangeTime() uint64 {
	return d.ChangeTime
}

func (d DirItem) GetAccessTime() uint64 {
	return d.AccessTime
}

func NewDirItem(fileOrDir IFileOrDir, name string, mode uint64, createTime uint64, modifyTime uint64, changeTime uint64, accessTime uint64) DirItem {
	return DirItem{fileOrDir.Hash(), name, mode, fileOrDir.Size(), fileOrDir.Count(), fileOrDir.TotalCount(), createTime, modifyTime, changeTime, accessTime}
}
