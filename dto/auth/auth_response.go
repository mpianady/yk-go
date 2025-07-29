package auth

type AuthResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGci..."`
	RefreshToken string `json:"refresh_token" example:"d8a53f1b..."`
}
