package models

import (
	"time"
)

type ProjectControl struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID   uint      `json:"project_id" binding:"required"`
	ControlID   uint      `json:"control_id" binding:"required"`
	CreatedBYID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Project       *Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedBYID"`
}

func (ProjectControl) TableName() string {
	return "project_controls"
}
