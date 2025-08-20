package models

import (
	"time"
)

// ComplianceAssessmentRequirementDocument represents a document attached to a compliance assessment requirement.
type ComplianceAssessmentRequirementDocumentTitle struct {
	ID                                               uint      `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceAssessmentRequirementResponsibleUnitID uint      `json:"compliance_assessment_requirement_responsible_unit_id" gorm:"not null"`
	RequiredDocumentTitleID                          *uint     `json:"required_document_title_id"`
	ApplicabilityStatusID                            uint      `json:"applicability_status_id"`
	CreatedBYID                                      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                                        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceAssessmentRequirementResponsibleUnit *ComplianceAssessmentRequirementResponsibleUnit    `gorm:"foreignKey:ComplianceAssessmentRequirementResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RequiredDocumentTitle                          *RequiredDocumentTitle                             `gorm:"foreignKey:RequiredDocumentTitleID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	Documents                                      *[]ComplianceAssessmentRequirementRequiredDocument `gorm:"foreignKey:ComplianceAssessmentRequirementDocumentTitleID"`
	CreatedByUser                                  *User                                              `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the default table name
func (ComplianceAssessmentRequirementDocumentTitle) TableName() string {
	return "compliance_assessment_requirement_document_titles"
}
