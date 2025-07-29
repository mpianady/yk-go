package auth
type RegisterResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR..."`
	RefreshToken string `json:"refresh_token" example:"dGhpc2lzYXJlZnJlc2h0b2tlbg=="`
}