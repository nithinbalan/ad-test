package models

import (
	"time"
)

type RiskSLA struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	SLARegisterID  uint      `json:"sla_register_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SLARegister   *SLARegister  `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (RiskSLA) TableName() string {
	return "risk_slas"
}
