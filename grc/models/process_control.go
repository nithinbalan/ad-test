package models

import (
	"time"
)

type ProcessControl struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID   uint      `json:"process_id" binding:"required"`
	ControlID   uint      `json:"control_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process       *Process `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedByID"`
}

func (ProcessControl) TableName() string {
	return "process_controls"
}
