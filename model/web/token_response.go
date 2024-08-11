package web

type TokenResponse struct {
	AccessToken       string `json:"access"`
	AccessValidUntil  string `json:"access_valid_until"`
	RefreshToken      string `json:"refresh"`
	RefreshValidUntil string `json:"refresh_valid_until"`
}
