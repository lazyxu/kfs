package dao

type IFileOrDir interface {
	Hash() string
	Size() uint64
	Count() uint64
	TotalCount() uint64
}

type FileOrDir struct {
	hash string
	size uint64
}

func (i FileOrDir) Hash() string {
	return i.hash
}

func (i FileOrDir) Size() uint64 {
	return i.size
}

func (i FileOrDir) Count() uint64 {
	return 1
}

func (i FileOrDir) TotalCount() uint64 {
	return 1
}
