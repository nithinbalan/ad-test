package models

import (
	"time"
)

// ErrorLog represents the schema of the 'error_logs' table with validation rules
type ErrorLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Level     string    `gorm:"size:45" validate:"required,max=45"`
	URL       string    `gorm:"size:45" validate:"max=45"`
	Message   string    `gorm:"type:varchar(MAX)" validate:"required"`
	CreatedAt time.Time `gorm:"type:datetime" validate:"required"`
}
