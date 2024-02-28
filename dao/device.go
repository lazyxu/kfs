package dao

type Device struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	OS        string `json:"os"`
	UserAgent string `json:"userAgent"`
	Hostname  string `json:"hostname"`
}
