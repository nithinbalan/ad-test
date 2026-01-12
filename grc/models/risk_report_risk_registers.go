package models

import "time"

type RiskReportRiskRegister struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID   uint      `json:"risk_report_id" gorm:"not null"`
	RiskRegisterID uint      `json:"risk_register_id" gorm:"not null"`
	CreatedBYID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskReport    *RiskReport   `gorm:"foreignKey:RiskReportID"`
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskReportRiskRegister) TableName() string {
	return "risk_report_risk_registers"
}
