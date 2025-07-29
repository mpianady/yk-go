package auth

import (
	"go-blog/models/user"
	"time"
)

type UserResponse struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewUserResponse(u user.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Role:        u.Role,
		Status:      u.Status,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}