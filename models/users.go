package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    FirstName string    `gorm:"column:first_name;type:varchar(100)" json:"first_name,omitempty"`
    LastName  string    `gorm:"column:last_name;type:varchar(100)" json:"last_name,omitempty"`
    Email     string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
    JoinedAt  time.Time `gorm:"joined_at" json:"joined_at"`
    UID       string    `gorm:"column:uid;uniqueIndex;not null" json:"uid"`
}
