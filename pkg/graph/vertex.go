package graph

import (
	"os"
)

type VertexWrapper interface {
	Hash() []byte
	IsLoaded() bool
	GetContent() Vertex
}

type Vertex interface {
	Hash() []byte

	IsDir() bool       // abbreviation for Mode().IsDir()
	Size() int64       // length in bytes for regular files; system-dependent for others
	Name() string      // base name of the file
	Mode() os.FileMode // file mode bits
	CreateTime() int64
	AccessTime() int64
	ModifyTime() int64
	ChangeTime() int64

	InList() []VertexWrapper
	OutList() []VertexWrapper
	EdgeList(typ string) []VertexWrapper
	IsSink() bool
	IsSource() bool
}
