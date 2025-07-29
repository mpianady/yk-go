package auth

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
