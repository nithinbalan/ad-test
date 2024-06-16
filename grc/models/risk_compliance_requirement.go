package models

import (
	"time"
)

type RiskComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID          uint      `json:"risk_register_id" binding:"required"`
	ComplianceRequirementID uint      `json:"compliance_requirement_id" binding:"required"`
	CreatedBy               uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegister          *RiskRegister          `gorm:"foreignkey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignkey:CreatedBy"`
}

func (RiskComplianceRequirement) TableName() string {
	return "risk_compliance_requirements"
}
