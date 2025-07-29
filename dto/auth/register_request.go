package auth

type RegisterRequest struct {
	Email     string `json:"email" example:"user@example.com" binding:"required,email"`
	Password  string `json:"password" example:"secret123" binding:"required,strong_password"`
	FirstName string `json:"first_name" example:"John" binding:"required"`
	LastName  string `json:"last_name" example:"Doe" binding:"required"`
}
