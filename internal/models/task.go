package models

import (
	"time"
)

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);default:'pending'"` // pending, in_progress, done
	DueDate     *time.Time
	UserID      uint `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
