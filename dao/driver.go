package dao

type Driver struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Typ         string `json:"type"`

	// sync
	Sync bool `json:"sync"`
	H    int  `json:"h"`
	M    int  `json:"m"`
	S    int  `json:"s"`

	// baidu photo
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`

	// local file
	DeviceId uint64 `json:"deviceId"`
	SrcPath  string `json:"srcPath"`
	Encoder  string `json:"encoder"`
}
