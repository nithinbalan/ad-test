package models

import "time"

type ComplianceRequirementResponsibleUnit struct {
	ID                                   uint       `json:"id" gorm:"primary_key;autoIncrement"`
	ComplianceRequirementID              uint       `json:"compliance_requirement_id" gorm:"not null"`
	ResponsibleUnitID                    uint       `json:"responsible_unit_id" gorm:"not null"`
	ComplianceAssessmentFrameworkValueID *uint      `json:"compliance_assessment_framework_value_id" gorm:"default:null"`
	CreatedByID                          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ComplianceRequirement              *ComplianceRequirement              `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUnit                    *OrganizationUnit                   `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceAssessmentFrameworkValue *ComplianceAssessmentFrameworkValue `gorm:"foreignKey:ComplianceAssessmentFrameworkValueID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                      *User                               `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                      *User                               `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the default table name
func (ComplianceRequirementResponsibleUnit) TableName() string {
	return "compliance_requirement_responsible_units"
}
