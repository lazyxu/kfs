package dao

type Driver struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Typ          string `json:"type"`
	Sync         bool   `json:"sync"`
	H            int    `json:"h"`
	M            int    `json:"m"`
	S            int    `json:"s"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
