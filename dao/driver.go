package dao

type Driver struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IDriver interface {
	GetName() string
	GetDescription() string
	GetCommitId() uint64
	GetSize() uint64
	GetCount() uint64
}

func (b Driver) GetName() string {
	return b.Name
}

func (b Driver) GetDescription() string {
	return b.Description
}

func NewDriver(name string, description string) Driver {
	return Driver{name, description}
}
