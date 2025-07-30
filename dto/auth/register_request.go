package auth

type RegisterRequest struct {
	Email     string `json:"email" example:"user@example.com" binding:"required,email" example:"newuser@example.com"`
	Password  string `json:"password" example:"secret123" binding:"required,strong_password" example:"strong@!password123"`
	FirstName string `json:"first_name" example:"John" binding:"required" example:"Alice"`
	LastName  string `json:"last_name" example:"Doe" binding:"required" example:"Dupont"`
}
