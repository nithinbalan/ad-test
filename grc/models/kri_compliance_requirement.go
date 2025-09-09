package models

import (
	"time"
)

type KriComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KriRegisterID           uint      `json:"kri_register_id" binding:"required"`
	ComplianceRegisterID    uint      `json:"compliance_register_id" binding:"required"`
	ComplianceRequirementID *uint     `json:"compliance_requirement_id" gorm:"default:null"`
	ComplianceClauseID      uint      `json:"compliance_clause_id" binding:"required"`
	CreatedByID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KRIRegister           *KRIRegister           `gorm:"foreignKey:KriRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister    *ComplianceRegister    `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause      *ComplianceClause      `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
}

func (KriComplianceRequirement) TableName() string {
	return "kri_compliance_requirements"
}
