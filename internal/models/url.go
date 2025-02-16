package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint `gorm:"primaryKey"`
	ShortCode string `gorm:"uniqueIndex;not null"`
	LongUrl   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	deletedAt gorm.DeletedAt `gorm:"index"`
}
