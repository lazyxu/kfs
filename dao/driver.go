package dao

type Driver struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Typ          string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
