package models

import (
	"time"
)

type AuditFindingComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID          uint      `json:"audit_finding_id" binding:"required"`
	ComplianceRegisterID    uint      `json:"compliance_register_id" binding:"required"`
	ComplianceRequirementID *uint     `json:"compliance_requirement_id" gorm:"default:null"`
	ComplianceClauseID      uint      `json:"compliance_clause_id" binding:"required"`
	CreatedByID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditFinding          *AuditFinding          `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister    *ComplianceRegister    `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause      *ComplianceClause      `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
}

func (AuditFindingComplianceRequirement) TableName() string {
	return "audit_finding_compliance_requirements"
}
