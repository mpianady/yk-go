package auth

import (
	"go-blog/models/user"
	"time"
)

type RefreshToken struct {
	ID uint `gorm:"primaryKey"`
	Token string `gorm:"not null;unique"`
	UserID uint `gorm:"not null;index"`
	User user.User `gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
