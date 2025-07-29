package post

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID         uint           `gorm:"primaryKey"`
	Title      string         `gorm:"not null"`
	Excerpt    string         `gorm:"not null"`
	Content    string         `gorm:"not null"`
	Categories []Category     `gorm:"many2many:post_categories;" json:"categories"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
