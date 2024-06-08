package models

import (
	"time"
)

type MeetingRisk struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID      uint      `json:"meeting_id" binding:"required"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	CreatedBy      uint      `json:"created_by" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Meeting       *Meeting      `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister  *RiskRegister `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignkey:CreatedBy"`
}

func (MeetingRisk) TableName() string {
	return "meeting_risks"
}
