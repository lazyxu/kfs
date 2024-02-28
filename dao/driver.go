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
	DeviceId string `json:"deviceId"`
	SrcPath  string `json:"srcPath"`
	Ignores  string `json:"ignores"`
	Encoder  string `json:"encoder"`
}

type DCIMDriver struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Typ         string `json:"type"`

	// DCIM
	MetadataList []Metadata `json:"metadataList"`
}

type DirCalculatedInfo struct {
	FileSize          uint64 `json:"fileSize"`
	FileCount         uint64 `json:"fileCount"`
	DistinctFileSize  uint64 `json:"distinctFileSize"`
	DistinctFileCount uint64 `json:"distinctFileCount"`
	DirCount          uint64 `json:"dirCount"`
}
