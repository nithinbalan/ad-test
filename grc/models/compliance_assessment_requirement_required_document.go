package models

import (
	"time"
)

type ComplianceAssessmentRequirementRequiredDocument struct {
	ID                                             uint       `json:"id" gorm:"primary_key;auto_increment"`
	ComplianceAssessmentRequirementDocumentTitleID uint       `json:"compliance_requirement_document_title_id"`
	Language                                       string     `json:"language"`
	URL                                            string     `json:"url"`
	Name                                           string     `json:"name" binding:"required"`
	Size                                           int        `json:"size"`
	Token                                          string     `json:"token"`
	OriginalName                                   string     `json:"original_name"`
	Type                                           string     `json:"type"`
	ApprovedByID                                   *uint      `json:"approved_by_id"`
	ApprovedAt                                     *time.Time `json:"approved_at"`
	RejectedByID                                   *uint      `json:"rejected_by_id"`
	RejectedAt                                     *time.Time `json:"rejected_at"`
	CreatedByID                                    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                                    uint       `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                                      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                                      time.Time  `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceAssessmentRequirementDocumentTitle *ComplianceAssessmentRequirementDocumentTitle `gorm:"foreignKey:ComplianceAssessmentRequirementDocumentTitleID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser                                *User                                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                                *User                                         `gorm:"foreignKey:UpdatedByID"`
	ApprovedByUser                               *User                                         `gorm:"foreignKey:ApprovedByID"`
	RejectedByUser                               *User                                         `gorm:"foreignKey:RejectedByID"`
}

// TableName overrides the table name used by ComplianceAssessmentRequirementRequiredDocument to `compliance_assessment_requirement_required_documents`
func (ComplianceAssessmentRequirementRequiredDocument) TableName() string {
	return "compliance_assessment_requirement_required_documents"
}
