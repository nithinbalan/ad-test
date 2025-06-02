package models

import "time"

type TaskSchedularLog struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	TaskName    string     `gorm:"column:task_name;not null"`
	Status      string     `gorm:"column:status;not null"`
	Message     *string    `gorm:"column:message"`
	RunAt       time.Time  `gorm:"column:run_at;not null;default:CURRENT_TIMESTAMP"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}
