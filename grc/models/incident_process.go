package models

import (
	"time"
)

type IncidentProcess struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	ProcessID          uint      `json:"process_id" binding:"required"`
	CreatedBy          uint      `json:"created_by" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignkey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process          *Process          `gorm:"foreignkey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignkey:CreatedBy"`
}

func (IncidentProcess) TableName() string {
	return "incident_processes"
}
