package object

import "os"

func (i *Metadata) Builder() *MetadataBuilder {
	return &MetadataBuilder{
		mode:       i.mode,
		birthTime:  i.birthTime,
		modifyTime: i.modifyTime,
		changeTime: i.changeTime,
		name:       i.name,
		size:       i.size,
		hash:       i.hash,
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

func (m *MetadataBuilder) Size(size int64) *MetadataBuilder {
	m.size = size
	return m
}

func (m *MetadataBuilder) Build() *Metadata {
	return &Metadata{
		mode:       m.mode,
		birthTime:  m.birthTime,
		modifyTime: m.modifyTime,
		changeTime: m.changeTime,
		name:       m.name,
		size:       m.size,
		hash:       m.hash,
	}
}
