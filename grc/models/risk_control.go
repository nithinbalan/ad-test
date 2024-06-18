package models

import (
	"time"
)

type RiskControl struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	ControlID      uint      `json:"control_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control      `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (RiskControl) TableName() string {
	return "risk_controls"
}
