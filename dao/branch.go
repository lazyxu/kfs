package dao

type Branch struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CommitId    uint64 `json:"commitId"`
	Size        uint64 `json:"size"`
	Count       uint64 `json:"count"`
}

type IBranch interface {
	GetName() string
	GetDescription() string
	GetCommitId() uint64
	GetSize() uint64
	GetCount() uint64
}

func (b Branch) GetName() string {
	return b.Name
}

func (b Branch) GetDescription() string {
	return b.Description
}

func (b Branch) GetCommitId() uint64 {
	return b.CommitId
}

func (b Branch) GetSize() uint64 {
	return b.Size
}

func (b Branch) GetCount() uint64 {
	return b.Count
}

func NewBranch(name string, commit Commit, dir Dir) Branch {
	return Branch{name, "", commit.Id, dir.size, dir.totalCount}
}
