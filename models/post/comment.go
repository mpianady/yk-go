package post

import (
	"go-blog/models/user"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey"`
	PostID    uint           `gorm:"index;not null"`
	Post      Post           `gorm:"foreignKey:PostID"`
	UserID    uint           `gorm:"index;not null"`
	User      user.User      `gorm:"foreignKey:UserID"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"-"`
	Children  []Comment      `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Status    string         `gorm:"type:ENUM('PENDING','APPROVED', 'REJECTED');default:'PENDING';not null"`
	Content   string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
