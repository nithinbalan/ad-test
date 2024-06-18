package models

import (
	"time"
)

type RiskProcess struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	ProcessID      uint      `json:"process_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process       *Process      `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (RiskProcess) TableName() string {
	return "risk_processes"
}
