package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"uniqueIndex;not null"`
	LongUrl   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId uint `gorm:"column:user_id;index"`
	User User `gorm:"foreignKey:UserId"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
