package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentRequirementDetailDTO struct {
	ID                     int                                                `json:"id"`
	ComplianceAssessment   *ComplianceAssessmentDetailDTO                     `json:"compliance_assessment"`
	ComplianceRequirement  ComplianceRequirementDetailDTO                     `json:"compliance_requirement"`
	IsCompleted            bool                                               `json:"is_completed"`
	Score                  *int                                               `json:"score"`
	Comments               *string                                            `json:"comments"`
	ComplianceStatusID     *uint                                              `json:"compliance_status_id"`
	ComplianceStatus       *ComplianceStatusDTO                               `json:"compliance_status"`
	CreatedAt              time.Time                                          `json:"created_at"`
	UpdatedAt              *time.Time                                         `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO                    `json:"created_by_user"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO                   `json:"updated_by_user"`
	RequiredDocumentTitles []*ComplianceAssessmentRequirementDocumentTitleDTO `json:"required_document_titles"`
}

type ComplianceAssessmentRequirementUpdateFormDTO struct {
	Score              int     `json:"score" binding:"required"`
	ComplianceStatusID uint    `json:"compliance_status_id" binding:"required"`
	Comments           *string `json:"comments"`
}

type ComplianceAssessmentRequirementDocumentTitleDTO struct {
	ID                                *uint                               `json:"id"`
	ComplianceAssessmentRequirementID uint                                `json:"compliance_assessment_requirement_id"`
	RequiredDocumentTitleID           *uint                               `json:"required_document_title_id"`
	RequiredDocumentTitle             masterDTOs.RequiredDocumentTitleDTO `json:"required_document_title"`
}

type ComplianceAssessmentRequirementDocumentFormDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
