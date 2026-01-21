package models

import (
	"time"
)

type ViolationReportComplianceRegister struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportID    uint      `json:"violation_report_id" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_register_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReport    *ViolationReport    `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (ViolationReportComplianceRegister) TableName() string {
	return "violation_report_compliance_registers"
}
