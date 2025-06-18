package models

import (
	"time"
)

type KPIRisk struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID  uint      `json:"kpi_register_id" binding:"required"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister   *KPIRegister  `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
}

func (KPIRisk) TableName() string {
	return "kpi_risks"
}
