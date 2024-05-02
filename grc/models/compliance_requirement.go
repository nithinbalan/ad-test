package models

import (
	"time"
)

type ComplianceRequirement struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceClauseID              uint       `json:"compliance_clause_id" binding:"required" gorm:"not null"`
	Title                           string     `json:"title" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	ComplianceRequirementCategoryID uint       `json:"compliance_requirement_category_id" binding:"required" gorm:"not null"`
	SeverityLevelID                 uint       `json:"severity_level_id" binding:"required" gorm:"not null"`
	ResponsibleUserID               uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Deadline                        *time.Time `json:"deadline" binding:"omitempty" gorm:"type:date"`
	ComplianceStatusID              uint       `json:"compliance_status_id" binding:"required" gorm:"not null"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy                       uint       `json:"created_by" gorm:"not null"`
	UpdatedBy                       *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	ComplianceClause              *ComplianceClause              `gorm:"foreignkey:ComplianceClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus              *ComplianceStatus              `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirementCategory *ComplianceRequirementCategory `gorm:"foreignkey:ComplianceRequirementCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel                 *SeverityLevel                 `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser               *User                          `gorm:"foreignkey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                 *User                          `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser                 *User                          `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceRequirement to `compliance_requirements`
func (ComplianceRequirement) TableName() string {
	return "compliance_requirements"
}
