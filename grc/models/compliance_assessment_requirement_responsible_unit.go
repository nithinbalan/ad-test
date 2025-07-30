package models

import "time"

type ComplianceAssessmentRequirementResponsibleUnit struct {
	ID                                      uint       `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceAssessmentRequirementID       uint       `json:"compliance_assessment_requirement_id" gorm:"not null"`
	ResponsibleUnitID                       uint       `json:"responsible_unit_id" gorm:"not null"`
	Deadline                                *time.Time `json:"deadline" gorm:"default:null"`
	ComplianceStatusID                      *uint      `json:"compliance_status_id" gorm:"default:null"`
	ComplianceAssessmentFrameworkValueID    *uint      `json:"compliance_assessment_framework_value_id" gorm:"default:null"`
	Score                                   *int       `json:"score" gorm:"default:null"`
	Comments                                *string    `json:"comments" gorm:"type:nvarchar(max);default:null"`
	IsCompleted                             bool       `json:"is_completed" gorm:"not null"`
	ComplianceAssessmentRequirementStatusID *uint      `json:"compliance_assessment_requirement_status_id" gorm:"default:null"`
	CreatedBYID                             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                             *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                               *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ComplianceAssessmentRequirement       *ComplianceAssessmentRequirement       `gorm:"foreignKey:ComplianceAssessmentRequirementID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	ResponsibleUnit                       *OrganizationUnit                      `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	ComplianceStatus                      *ComplianceStatus                      `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	ComplianceAssessmentFrameworkValue    *ComplianceAssessmentFrameworkValue    `gorm:"foreignKey:ComplianceAssessmentFrameworkValueID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatus `gorm:"foreignKey:ComplianceAssessmentRequirementStatusID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser                         *User                                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	UpdatedByUser                         *User                                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

func (ComplianceAssessmentRequirementResponsibleUnit) TableName() string {
	return "compliance_assessment_requirement_responsible_units"
}
