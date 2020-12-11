package object

import (
	"os"
	"time"
)

type Metadata struct {
	Mmode      os.FileMode
	BirthTime  int64
	ModifyTime int64
	ChangeTime int64
	Size       int64
	Name       string
	Hash       string
}

func (m *Metadata) Builder() *MetadataBuilder {
	return &MetadataBuilder{
		mode:       m.Mmode,
		birthTime:  m.BirthTime,
		modifyTime: m.ModifyTime,
		changeTime: m.ChangeTime,
		name:       m.Name,
		size:       m.Size,
		hash:       m.Hash,
	}
}

type MetadataBuilder struct {
	mode       os.FileMode
	birthTime  int64
	modifyTime int64
	changeTime int64
	name       string
	size       int64
	hash       string
}

func (m *MetadataBuilder) Mode(mode os.FileMode) *MetadataBuilder {
	m.mode = mode
	return m
}

func (m *MetadataBuilder) Hash(hash string) *MetadataBuilder {
	m.hash = hash
	return m
}

func (m *MetadataBuilder) Name(name string) *MetadataBuilder {
	m.name = name
	return m
}

func (m *MetadataBuilder) ChangeTime(changeTime int64) *MetadataBuilder {
	m.changeTime = changeTime
	return m
}

func (m *MetadataBuilder) ModifyTime(modifyTime int64) *MetadataBuilder {
	m.modifyTime = modifyTime
	return m
}

func (m *MetadataBuilder) Build() *Metadata {
	return &Metadata{
		Mmode:      m.mode,
		BirthTime:  m.birthTime,
		ModifyTime: m.modifyTime,
		ChangeTime: m.changeTime,
		Name:       m.name,
		Size:       m.size,
		Hash:       m.hash,
	}
}

func (i *Metadata) Mode() os.FileMode {
	return i.Mmode
}

func (i *Metadata) IsFile() bool {
	return i.Mmode&S_IFREG != 0
}

func (i *Metadata) IsDir() bool {
	return i.Mmode&S_IFDIR != 0
}

func (base *Obj) NewDirMetadata(name string, perm os.FileMode) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		Mmode:      S_IFDIR | (perm & os.ModePerm),
		BirthTime:  now,
		ModifyTime: now,
		ChangeTime: now,
		Name:       name,
		Size:       0,
		Hash:       base.EmptyDirHash,
	}
}

func (base *Obj) NewFileMetadata(name string, perm os.FileMode) *Metadata {
	now := time.Now().UnixNano()
	return &Metadata{
		Mmode:      S_IFREG | (perm & os.ModePerm),
		BirthTime:  now,
		ModifyTime: now,
		ChangeTime: now,
		Name:       name,
		Size:       0,
		Hash:       base.EmptyFileHash,
	}
}
