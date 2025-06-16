package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Nickname  string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // захэшированный пароль
	CreatedAt time.Time
	UpdatedAt time.Time
	Tasks     []Task `gorm:"foreignKey:UserID"`
}
