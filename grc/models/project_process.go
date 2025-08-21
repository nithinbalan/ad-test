package models

import (
	"time"
)

type ProjectProcess struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID   uint      `json:"process_id" binding:"required"`
	ProjectID   uint      `json:"project_id" binding:"required"`
	CreatedByID uint      `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process       *Process `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Project       *Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedByID"`
}

func (ProjectProcess) TableName() string {
	return "project_processes"
}
