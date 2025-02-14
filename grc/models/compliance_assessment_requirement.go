package models

import "time"

type ComplianceAssessmentRequirement struct {
	ID                                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceAssessmentID                  uint       `json:"compliance_assessment_id" binding:"required" gorm:"column:compliance_assessment_id;not null"`
	ComplianceRequirementID                 uint       `json:"compliance_requirement_id" binding:"required" gorm:"column:compliance_requirement_id;not null"`
	ComplianceStatusID                      *uint      `json:"compliance_status_id"`
	ComplianceAssessmentRequirementStatusID *uint      `json:"compliance_assessment_requirement_status_id" gorm:"default:null"`
	Comments                                *string    `json:"comments"`
	Score                                   *int       `json:"score"`
	IsCompleted                             bool       `json:"is_completed"`
	CreatedAt                               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                             *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceAssessment                  *ComplianceAssessment                           `gorm:"foreignKey:ComplianceAssessmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirement                 *ComplianceRequirement                          `gorm:"foreignKey:ComplianceRequirementID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus                      *ComplianceStatus                               `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatus          `gorm:"foreignKey:ComplianceAssessmentRequirementStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                         *User                                           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                         *User                                           `gorm:"foreignKey:UpdatedBYID"`
	RequiredDocumentTitles                []*ComplianceAssessmentRequirementDocumentTitle `gorm:"foreignKey:ComplianceAssessmentRequirementID"`
}

// TableName overrides the table name used by Compliance Assessment Requirement to `compliance_assessment_requirements`
func (ComplianceAssessmentRequirement) TableName() string {
	return "compliance_assessment_requirements"
}
