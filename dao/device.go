package dao

type Device struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	OS   string `json:"os"`
}
