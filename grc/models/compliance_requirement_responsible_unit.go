package models

import "time"

type ComplianceRequirementResponsibleUnit struct {
	ID                                   uint       `json:"id" gorm:"primary_key;autoIncrement"`
	ComplianceRequirementID              uint       `json:"compliance_requirement_id" gorm:"not null"`
	ResponsibleUnitID                    uint       `json:"responsible_unit_id" gorm:"not null"`
	ComplianceAssessmentFrameworkValueID *uint      `json:"compliance_assessment_framework_value_id" gorm:"default:null"`
	CompletedByUserID                    *uint      `json:"completed_by_user_id"`
	CompletedAt                          *time.Time `json:"completed_at" gorm:"type:datetime;default:null"`
	Comments                             *string    `json:"comments" gorm:"type:nvarchar(max);default:null"`
	CreatedByID                          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ComplianceRequirement              *ComplianceRequirement              `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUnit                    *OrganizationUnit                   `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceAssessmentFrameworkValue *ComplianceAssessmentFrameworkValue `gorm:"foreignKey:ComplianceAssessmentFrameworkValueID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CompletedByUser                    *User                               `gorm:"foreignKey:CompletedByUserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser                      *User                               `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                      *User                               `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the default table name
func (ComplianceRequirementResponsibleUnit) TableName() string {
	return "compliance_requirement_responsible_units"
}
