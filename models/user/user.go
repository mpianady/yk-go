package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint           `gorm:"primaryKey"`
	Email       string         `gorm:"unique"`
	Password    string         `gorm:"not null"`
	Role        string         `gorm:"type:ENUM('AUTHOR','CONTRIBUTOR','ADMIN','READER');default:'READER';not null"`
	Status      string         `gorm:"type:ENUM('ACTIVE','INACTIVE','BANNED','PENDING');default:'ACTIVE';not null"`
	FirstName   string         `gorm:"not null"`
	LastName    string         `gorm:"not null"`
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
