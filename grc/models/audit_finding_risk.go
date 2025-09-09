package models

import (
	"time"
)

type AuditFindingRisk struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID uint      `json:"audit_finding_id" binding:"required"`
	RiskRegisterID uint      `json:"risk_register_id" binding:"required"`
	CreatedBYID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	AuditFinding  *AuditFinding `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedBYID"`
}

func (AuditFindingRisk) TableName() string {
	return "audit_finding_risks"
}
