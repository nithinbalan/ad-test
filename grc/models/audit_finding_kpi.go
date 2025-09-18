package models

import (
	"time"
)

type AuditFindingKpi struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID  uint      `json:"kpi_register_id" binding:"required"`
	AuditFindingID uint      `json:"audit_finding_id" binding:"required"`
	CreatedByID    uint      `json:"created_by_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister   *KPIRegister `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Meeting       *Meeting     `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (AuditFindingKpi) TableName() string {
	return "audit_finding_kpis"
}
