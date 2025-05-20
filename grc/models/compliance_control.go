package models

import (
	"time"
)

type ComplianceControl struct {
	ID                      uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID    uint      `json:"compliance_register_id" binding:"required"`
	ComplianceRequirementID *uint     `json:"compliance_requirement_id" gorm:"default:null"`
	ControlID               uint      `json:"control_id" binding:"required"`
	ComplianceClauseID      *uint     `json:"compliance_clause_id" gorm:"default:null"`
	CreatedByID             uint      `json:"created_by_id" binding:"required"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister    *ComplianceRegister    `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control               *Control               `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause      *ComplianceClause      `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedByID"`
}

func (ComplianceControl) TableName() string {
	return "compliance_controls"
}
