package graph

type Graph interface {
	Sources() []VertexWrapper
	Version() int
}
