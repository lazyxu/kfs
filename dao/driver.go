package dao

type Driver struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Typ         string `json:"type"`

	// sync
	Sync bool  `json:"sync"`
	H    int64 `json:"h"`
	M    int64 `json:"m"`

	// baidu photo
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`

	// local file
	DeviceId uint64 `json:"deviceId"`
	SrcPath  string `json:"srcPath"`
	Ignores  string `json:"ignores"`
	Encoder  string `json:"encoder"`
}
