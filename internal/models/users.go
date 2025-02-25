package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	FirstName string    `gorm:"column:first_name;type:varchar(100)"`
	LastName  string    `gorm:"column:last_name;type:varchar(100)"`
	JoinedAt  time.Time `gorm:"column:uid;uniqueIndex;not null"`
	UID       string    `gorm:"index"`
}


