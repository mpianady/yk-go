package post

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string         `gorm:"not null"`
	Posts       []Post         `gorm:"many2many:post_categories;" json:"posts"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
