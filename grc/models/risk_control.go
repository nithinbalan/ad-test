package models

import (
	"time"
)

type RiskControl struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	ControlID      uint      `json:"control_id" binding:"required"`
	CreatedBy      uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control      `gorm:"foreignkey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignkey:CreatedBy"`
}

func (RiskControl) TableName() string {
	return "risk_controls"
}
