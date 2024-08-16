package models

import (
	"time"
)

type MeetingIncident struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID          uint      `json:"meeting_id" binding:"required"`
	IncidentRegisterID uint      `json:"risk_register_id" binding:"required"`
	CreatedBYID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Meeting          *Meeting          `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister *IncidentRegister `gorm:"foreignkey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignkey:CreatedBYID"`
}

func (MeetingIncident) TableName() string {
	return "meeting_incidents"
}
