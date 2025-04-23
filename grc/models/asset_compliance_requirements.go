package models

import (
	"time"
)

type AssetComplianceRequirement struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID         uint      `json:"asset_register_id" gorm:"not null"`
	ComplianceClauseID      *uint     `json:"compliance_clause_id"`
	ComplianceRequirementID *uint     `json:"compliance_requirement_id"`
	CreatedBYID             uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AssetRegister         *AssetRegister         `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceClause      *ComplianceClause      `gorm:"foreignKey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name used by GORM
func (AssetComplianceRequirement) TableName() string {
	return "asset_compliance_requirements"
}
