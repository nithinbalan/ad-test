package models

import (
	"time"
)

type RiskProcess struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	ProcessID      uint      `json:"process_id" binding:"required"`
	CreatedBy      uint      `json:"created_by" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process       *Process      `gorm:"foreignkey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignkey:CreatedBy"`
}

func (RiskProcess) TableName() string {
	return "risk_processes"
}
