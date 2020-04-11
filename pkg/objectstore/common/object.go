package common

import (
	"encoding/binary"
)

const (
	MagicNumber  = 0x005346616c616f4b
	MajorVersion = 0x00
	MinorVersion = 0x00
	PatchVersion = 0x00
)

const (
	TypeBlob   = 0x00000000
	TypeTree   = 0x00000001
	TypeCommit = 0x00000002
)

type Header struct {
	MagicNumber uint64
	Version     uint32
	ObjectType  uint32
	DataSize    uint64
	NewLine     uint64
}

func version(endian binary.ByteOrder) uint32 {
	return endian.Uint32([]byte{MajorVersion, MinorVersion, PatchVersion, 0x00})
}

func NewHeader(typ uint32, size uint64, endian binary.ByteOrder) *Header {
	return &Header{
		MagicNumber: MagicNumber,
		Version:     version(endian),
		ObjectType:  typ,
		DataSize:    size,
		NewLine:     0x0a00000000000000,
	}
}
