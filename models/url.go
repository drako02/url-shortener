package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ShortCode string         `gorm:"uniqueIndex;not null" json:"short_code"`
	LongUrl   string         `gorm:"not null" json:"long_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UserId    uint           `gorm:"column:user_id;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
